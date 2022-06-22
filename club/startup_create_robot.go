//机器人创建
package club

import (
	"fmt"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	//global control variable
	capacity          int           = 4000
	taskGenerateCycle time.Duration = time.Second * 300 //新任务监听周期
	taskCrontabCycle  time.Duration = time.Second * 300 //定时任务检测周期
	taskProcessNum    int           = 10                //任务处理协程数量

	taskChannel chan int64      = make(chan int64, capacity)
	taskMapping map[int64]int64 = make(map[int64]int64, capacity)
	mu          *sync.Mutex     = new(sync.Mutex)
)

func StartupGuildCreateRobot() {
	// read from db -----> write to taskMapping
	go taskGenerate()
	// read from taskMapping ----->  write to taskChannel
	go taskCrontab()

	// consume from taskChannel
	for i := 0; i < taskProcessNum; i++ {
		go taskProcess()
	}
}

func taskGenerate() {
	logPrefix := "taskGenerate: "

	for {
		var logMsg string
		start := time.Now()

		//usage like do{}while()
		for {
			//sacn all guild
			list, err := service.GetAllGuildInfoFromDB()
			if err != nil {
				time.Sleep(time.Second * 10)
				logMsg = "GetAllGuildInfoFromDB " + err.Error()
				break
			}

			logMsg = logMsg + fmt.Sprintf("Guild length:%d", len(list))

			//deleteJob
			var deleteJob []model.Guild
			var newJob []model.Guild

			for _, v := range list {
				if v.DeletedAt.Valid == true {
					deleteJob = append(deleteJob, v)
					continue
				}
			}

			//deletejob proc
			if len(deleteJob) > 0 {
				mu.Lock()
				for _, v := range deleteJob {
					delete(taskMapping, v.ID)
				}
				mu.Unlock()
			}

			//newjob proc
			if len(newJob) > 0 {
				mu.Lock()
				for _, v := range newJob {
					taskMapping[v.ID] = getActiveTimeByRand()
				}
				mu.Unlock()
			}

			logMsg = logMsg + " " + fmt.Sprintf("deleteJob:%d, newJob:%d", len(deleteJob), len(newJob))
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club").Info(logPrefix, logMsg, "cost:", cost, "ms")

		time.Sleep(taskGenerateCycle)
	}
}

func taskCrontab() {
	logPrefix := "taskCrontab: "
	for {

		start := time.Now()
		now := start.Unix()

		var needProcess []int64

		mu.Lock()
		for guildID, activeTime := range taskMapping {
			if now-activeTime >= 0 {
				needProcess = append(needProcess, guildID)
				//update next active time
				taskMapping[guildID] = getActiveTimeByRand()
			}
		}
		mu.Unlock()

		for _, activeTaskGuilID := range needProcess {
			taskChannel <- activeTaskGuilID
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club").Info(logPrefix, "need to process num:", len(needProcess), "cost:", cost, "ms")

		time.Sleep(taskCrontabCycle)
	}
}

func taskProcess() {
	logPrefix := "taskProcess: "
	for {

		guilID := <-taskChannel

		start := time.Now()

		logSuffix := fmt.Sprintf("guild_id:%d", guilID)

		var logMsg string
		for {
			//1.从DB获取工会信息
			guildInfo, err := service.GetGuildInfoByGuildID(guilID)
			if err != nil {
				logMsg = "GetGuildInfoByGuildID error: " + err.Error()
				break
			}

			//2.获取机器人和真实用户分布
			userDistributions, err := service.GetGuildUserTypeDistribution(guilID)
			if err != nil {
				logMsg = "service.GetGuildUserTypeDistribution error: " + err.Error()
				break
			}

			var robotUsers []model.User
			var normalUsers []model.User
			if v, ok := userDistributions[service.USERTYPE_CLUB_ROBOT_SERVICE]; ok {
				robotUsers = v
			}
			if v, ok := userDistributions[service.USERTYPE_NORMAL]; ok {
				normalUsers = v
			}

			//3.判断是否满足创建机器人的条件
			reason, ok := IsTheGuildCanUsingRobot(guildInfo, robotUsers, normalUsers)
			if !ok {
				logMsg = reason
				break
			}

			//4.创建机器人
			//4.1获取随机创建机器人数量
			newNum := getGenerateRobotNumByRand()
			robotMaxLimitNum := getRobotMaxLimitNum()
			//4.2边界检查
			if len(robotUsers)+newNum > robotMaxLimitNum {
				newNum = robotMaxLimitNum - len(robotUsers)
			}

			//newNum一定>0

			//4.3创建机器人
			var newUids []int64
			isCreateOk := true
			for i := 0; i < newNum; i++ {
				if rbtUser, err := CreateRobot(guildInfo, robotUsers, normalUsers); err != nil {
					logMsg = "CreateRobot error: " + err.Error()
					isCreateOk = false
					break
				} else {
					//名称过滤用
					robotUsers = append(robotUsers, rbtUser)
					//日志打印和加入工会用
					newUids = append(newUids, rbtUser.UserID)
				}
			}
			if !isCreateOk {
				//只要有一个失败，则不执行加入工会操作
				break
			}
			//5.加入工会
			for _, uid := range newUids {
				if err := service.JoinToGuild(guilID, uid); err != nil {
					logMsg = logMsg + " " + "JoinToGuild error:" + err.Error() + " " + fmt.Sprintf("uid:%d", uid)
					break
				} else {
					//加入执行第一次任务
					addFirstAction(guilID, uid)
				}
			}

			logMsg = fmt.Sprintf("CreateRobot newNum:%d, createOkNum:%d [%v] Join club success", newNum, len(newUids), newUids)
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club").Info(logPrefix, logMsg, logSuffix, "cost:", cost, "ms")
	}
}
