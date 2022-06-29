package create

import (
	"fmt"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/club/action"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	taskChannel chan int64      = make(chan int64, 2000)
	taskMapping map[int64]int64 = make(map[int64]int64, 5000)
	taskMu      *sync.Mutex     = new(sync.Mutex)
)

func Startup() {
	go taskUpdate(20)  //read from db guild -> write to taskMapping
	go taskCrontab(10) //read from taskMapping -> write to taskChannel
	go taskProcess()   //read from taskChannel -> operation create robot
}

func taskUpdate(t int) {
	logPrefix := "taskUpdate: "
	for {
		start := time.Now()
		logMsg := ""
		//usage like do{}while()
		for {
			//get all club
			list, err := service.GetAllGuildInfoFromDB()
			if err != nil {
				logMsg = "GetAllGuildInfoFromDB Error: " + err.Error()
				break
			}

			var (
				deleteJob []model.Guild
				newJob    []model.Guild
			)

			for _, v := range list {
				if v.DeletedAt.Valid == true {
					deleteJob = append(deleteJob, v)
				} else {
					newJob = append(newJob, v)
				}
			}

			//delete job proc
			if len(deleteJob) > 0 {
				taskMu.Lock()
				for _, v := range deleteJob {
					delete(taskMapping, v.ID)
				}
				taskMu.Unlock()
			}

			//new job proc
			if len(newJob) > 0 {
				taskMu.Lock()
				for _, v := range newJob {
					taskMapping[v.ID] = time.Now().Unix() + config.GetActiveTimeByRand()
				}
				taskMu.Unlock()
			}

			logMsg = logMsg + " " + fmt.Sprintf("deleteJob:%d, newJob:%d", len(deleteJob), len(newJob))
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-create").Info(logPrefix, logMsg, "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func taskCrontab(t int) {
	logPrefix := "taskCrontab: "
	for {

		start := time.Now()
		now := start.Unix()

		var needProcess []int64

		taskMu.Lock()
		for guildID, activeTime := range taskMapping {
			if now-activeTime >= 0 || guildID == 125323777000079360 {
				needProcess = append(needProcess, guildID)
				taskMapping[guildID] = config.GetActiveTimeByRand()
			}
		}
		total := len(taskMapping)
		taskMu.Unlock()

		for _, activeTaskGuilID := range needProcess {
			taskChannel <- activeTaskGuilID
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-create").Info(logPrefix, fmt.Sprintf("total:%d, needProcess:%d", total, len(needProcess)), "cost:", cost, "ms")

		time.Sleep(time.Duration(t) * time.Second)
	}
}

func taskProcess() {
	logPrefix := "taskProcess: "
	for {

		guilID := <-taskChannel

		start := time.Now()

		var logMsg string

		logSuffix := fmt.Sprintf("gid:%d", guilID)

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
			newNum := config.GetGenerateRobotNumByRand()
			robotMaxLimitNum := config.GetRobotMaxLimitNum()
			//4.2边界检查
			if len(robotUsers)+newNum > robotMaxLimitNum {
				newNum = robotMaxLimitNum - len(robotUsers)
			}

			//newNum一定>0

			//4.3创建机器人
			var newRobots []model.User
			var newRobotsUid []int64
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
					newRobots = append(newRobots, rbtUser)
					newRobotsUid = append(newRobotsUid, rbtUser.UserID)
				}
			}
			if !isCreateOk {
				//只要有一个失败，则不执行加入工会操作
				break
			}

			//5.加入工会
			for _, ru := range newRobots {
				if _, err := service.SendJoinToGuildRPC(ru.AccountID, ru.UserID, guilID); err != nil {
					logMsg = fmt.Sprintf("robotUser:%d, SendJoinToGuildRPC Error: %s", ru.UserID, err.Error())
					break
				}
				action.CreateFirstInJob(ru.UserID, guilID)
			}

			logMsg = fmt.Sprintf("CreateRobot targetNum:%d, createOkNum:%d %v, Join club success", newNum, len(newRobots), newRobotsUid)
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-create").Info(logPrefix, logMsg, logSuffix, "cost:", cost, "ms")
	}
}