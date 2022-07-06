package action

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/casual-server-lib/util"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/qa"
	"github.com/joycastle/matching-story-robot-service/service"
)

const (
	JOB_TYPE_FIRSTIN      = "FirstJoinGuild"
	JOB_TYPE_REQUEST      = "Request"
	JOB_TYPE_REQUEST_CHAT = "RequestChat"
	JOB_TYPE_HELP         = "Help"
	JOB_TYPE_ACTIVITY     = "Activity"

	JOB_TYPE_DELETE_GUILD = "DeleteGuild"

	JOB_TYPE_OWN_AI = "OwnAI"
)

type Job struct {
	GuildID    int64
	UserID     int64
	ActionTime int64
	RobotNum   int
}

var (
	capacityMap     int = 10000
	capacityChannel int = 2000

	//firstIn job other operation see action_first_in.go
	firstInCrontabJob     map[string]*Job = make(map[string]*Job, 1000)
	firstInCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	firstInProcessChannel chan *Job       = make(chan *Job, 1000)

	//request job other operation see action_request.go
	requestCrontabJob     map[string]*Job = make(map[string]*Job, capacityMap)
	requestCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	requestProcessChannel chan *Job       = make(chan *Job, capacityChannel)

	//request chat job other operation see action_request.go
	requestChatCrontabJob     map[string]*Job = make(map[string]*Job, 1000)
	requestChatCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	requestChatProcessChannel chan *Job       = make(chan *Job, 1000)

	//help job other operation see action_help.go
	helpCrontabJob     map[string]*Job = make(map[string]*Job, capacityMap)
	helpCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	helpProcessChannel chan *Job       = make(chan *Job, capacityChannel)

	//delete guild
	deleteGuildChannel chan *Job = make(chan *Job, 100)

	//using for ownaction
	ownActionCrontabJob     map[string]*Job = make(map[string]*Job, capacityMap)
	ownActionCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	ownActionProcessChannel chan *Job       = make(chan *Job, capacityChannel)

	//data source for update robot action
	RobotActionUpdateCrontabJob   map[string]*Job = make(map[string]*Job, capacityMap)
	RobotActionUpdateCrontabJobMu *sync.Mutex     = new(sync.Mutex)
)

func DeleteJob(targets map[string]*Job, mu *sync.Mutex, newTargets map[string]*Job) int {
	c := 0
	mu.Lock()
	for k, _ := range targets {
		if _, ok := newTargets[k]; !ok {
			delete(targets, k)
			c = c + 1
		}
	}
	mu.Unlock()
	return c
}

func CreateJob(targets map[string]*Job, mu *sync.Mutex, newTargets map[string]*Job, actionTimeHandler func() int64) int {
	c := 0
	mu.Lock()
	for k, newJob := range newTargets {
		if _, ok := targets[k]; !ok {
			if actionTimeHandler != nil {
				newJob.ActionTime = actionTimeHandler()
			} else {
				newJob.ActionTime = defaultActiveTimeHandler()
			}
			targets[k] = newJob
			c = c + 1
		}
	}
	mu.Unlock()
	return c
}

func JobKey(userID, guildID int64) string {
	return fmt.Sprintf("%d-%d", userID, guildID)
}

func CrontabGenerateJob(jobType string, targets map[string]*Job, mu *sync.Mutex, output chan *Job, cycleTimeHandler func(*Job) (int64, error), t int) {
	logPrefix := "CrontabGenerateJob-" + jobType
	for {
		start := time.Now()
		now := start.Unix()

		needProcess := []*Job{}
		needResetProcess := []*Job{}

		mu.Lock()
		for k, job := range targets {
			if now-job.ActionTime >= 0 || k == "130714000009-125323777000079360" {
				needProcess = append(needProcess, job)
				delete(targets, k)
			}
		}
		count := len(targets)
		mu.Unlock()

		needProcessNum := len(needProcess)
		if needProcessNum > 0 {
			//循环定时器
			if cycleTimeHandler != nil {
				//单独处理防止加锁时间过长
				for _, job := range needProcess {
					next, err := cycleTimeHandler(job)
					fmt.Println(logPrefix, job, util.FromUnixtime(job.ActionTime).Format("2006-01-02 15:04:05"), util.FromUnixtime(next).Format("2006-01-02 15:04:05"), err)
					if err != nil {
						log.Get("club-dispatch").Fatal(logPrefix, "ActiveTime set error using default ActiveTime:", err)
						next = defaultActiveTimeHandler()
					}

					job.ActionTime = next
					needResetProcess = append(needResetProcess, job)
				}
				//重新设置时间
				if len(needResetProcess) > 0 {
					mu.Lock()
					for _, job := range needProcess {
						k := JobKey(job.UserID, job.GuildID)
						targets[k] = job
					}
					mu.Unlock()
				}
			}

			for _, job := range needProcess {
				output <- job
			}
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-dispatch").Info(logPrefix, "total:", count, "needNum:", needProcessNum, "resetNun:", len(needResetProcess), "cost:", cost, "ms")

		time.Sleep(time.Duration(t) * time.Second)
	}
}

func Startup() {
	go UpdateRobotJobs(10)

	go CrontabGenerateJob(JOB_TYPE_FIRSTIN, firstInCrontabJob, firstInCrontabJobMu, firstInProcessChannel, nil, 10)
	go CrontabGenerateJob(JOB_TYPE_REQUEST, requestCrontabJob, requestCrontabJobMu, requestProcessChannel, nil, 10)
	go CrontabGenerateJob(JOB_TYPE_REQUEST_CHAT, requestChatCrontabJob, requestChatCrontabJobMu, requestChatProcessChannel, nil, 10)
	go CrontabGenerateJob(JOB_TYPE_HELP, helpCrontabJob, helpCrontabJobMu, helpProcessChannel, nil, 10)
	go CrontabGenerateJob(JOB_TYPE_OWN_AI, ownActionCrontabJob, ownActionCrontabJobMu, ownActionProcessChannel, cycleTimeHandlerOwnAi, 10)

	go JobActionProcess(JOB_TYPE_FIRSTIN, firstInProcessChannel, firstInActionHandler, false)
	go JobActionProcess(JOB_TYPE_REQUEST, requestProcessChannel, requestActionHandler, true)
	go JobActionProcess(JOB_TYPE_REQUEST_CHAT, requestChatProcessChannel, requestChatActionHandler, true)
	go JobActionProcess(JOB_TYPE_HELP, helpProcessChannel, helpActionHandler, true)
	go JobActionProcess(JOB_TYPE_OWN_AI, ownActionProcessChannel, ownActionHandler, true)

	go JobActionProcess(JOB_TYPE_DELETE_GUILD, deleteGuildChannel, deleteActionHandler, false)

	//更新配置周一0点
	go UpdateRobotConfigMonday(RobotActionUpdateCrontabJob, RobotActionUpdateCrontabJobMu)
}

func UpdateRobotJobs(t int) error {
	logPrefix := "UpdateRobotJobs: "
	for {
		start := time.Now()
		logMsg := ""

		for {

			userMap, err := service.GetAllGuildUserMapInfos()
			if err != nil {
				logMsg = "GetAllGuildUserMapInfo Error: " + err.Error()
				break
			}

			//是否开启QA调试功能
			if qa.OpenQaDebug() {
				debugIdsMap := qa.GetGuildIDMap()
				targetUserMap := []model.GuildUserMap{}
				for _, v := range userMap {
					if _, ok := debugIdsMap[v.GuildID]; ok {
						targetUserMap = append(targetUserMap, v)
					}
				}
				userMap = targetUserMap
			}

			if len(userMap) == 0 {
				logMsg = "no data to process"
				break
			}

			//获取uid和工会维度数据
			uids := []int64{}
			guildIDMaps := make(map[int64]struct{})
			guildIDUserCountMaps := make(map[int64]int)
			for _, v := range userMap {
				uids = append(uids, v.UserID)
				guildIDMaps[v.GuildID] = struct{}{}
				if _, ok := guildIDUserCountMaps[v.GuildID]; !ok {
					guildIDUserCountMaps[v.GuildID] = 0
				}
				guildIDUserCountMaps[v.GuildID]++
			}

			guildIDs := []int64{}
			for k, _ := range guildIDMaps {
				guildIDs = append(guildIDs, k)
			}

			//获取用户类型数据
			userTypes, err := service.GetUserInfosWithField(uids, []string{"user_type"})
			if err != nil {
				logMsg = "GetUserTypes Error: " + err.Error()
				break
			}

			filterUserTypeMaps := make(map[int64]string, len(userTypes))
			for _, v := range userTypes {
				filterUserTypeMaps[v.UserID] = v.UserType
			}

			//过滤工会状态
			guildDeleteInfos, err := service.GetGuildInfosWithField(guildIDs, []string{"deleted_at"})
			if err != nil {
				logMsg = "GetGuildInfosWithField Error: " + err.Error()
				break
			}
			filterGuildDeleteMap := make(map[int64]bool, len(guildIDs))
			for _, v := range guildDeleteInfos {
				filterGuildDeleteMap[v.ID] = v.DeletedAt.Valid
			}

			//判断是否全部是机器人
			guildRobotUserCountMaps := make(map[int64]int)
			for _, v := range userMap {
				ut, ok := filterUserTypeMaps[v.UserID]
				if !ok {
					continue
				}
				if ut == service.USERTYPE_CLUB_ROBOT_SERVICE {
					if _, ok := guildRobotUserCountMaps[v.GuildID]; !ok {
						guildRobotUserCountMaps[v.GuildID] = 0
					}
					guildRobotUserCountMaps[v.GuildID]++
				}
			}
			//获取需要删除的任务
			filterNeedDeleteGuildMaps := make(map[int64]bool)
			for k, allNum := range guildIDUserCountMaps {
				if robotNum, ok := guildRobotUserCountMaps[k]; ok && allNum == robotNum {
					if isDel, okk := filterGuildDeleteMap[k]; okk && !isDel {
						filterNeedDeleteGuildMaps[k] = true
						deleteGuildChannel <- &Job{GuildID: k, RobotNum: robotNum}
					}
				}
			}

			//获取可用robot
			robotUsers := []model.GuildUserMap{}
			for _, v := range userMap {
				if n, ok := filterNeedDeleteGuildMaps[v.GuildID]; ok && n {
					continue
				}

				ut, ok := filterUserTypeMaps[v.UserID]
				if !ok {
					continue
				}
				if ut != service.USERTYPE_CLUB_ROBOT_SERVICE {
					continue
				}

				isDel, ok := filterGuildDeleteMap[v.GuildID]
				if !ok {
					continue
				}
				if isDel {
					continue
				}

				robotUsers = append(robotUsers, v)
			}

			robotNewJobsMap := make(map[string]*Job, len(robotUsers))
			for _, v := range robotUsers {
				index := JobKey(v.UserID, v.GuildID)
				robotNewJobsMap[index] = &Job{
					GuildID: v.GuildID,
					UserID:  v.UserID,
				}
			}

			//delete job
			DeleteJob(firstInCrontabJob, firstInCrontabJobMu, robotNewJobsMap)
			DeleteJob(requestCrontabJob, requestCrontabJobMu, robotNewJobsMap)
			DeleteJob(requestChatCrontabJob, requestChatCrontabJobMu, robotNewJobsMap)
			DeleteJob(helpCrontabJob, helpCrontabJobMu, robotNewJobsMap)
			DeleteJob(ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap)
			DeleteJob(RobotActionUpdateCrontabJob, RobotActionUpdateCrontabJobMu, robotNewJobsMap)

			//create job
			CreateJob(requestCrontabJob, requestCrontabJobMu, robotNewJobsMap, requestActiveTimeHandler)
			CreateJob(helpCrontabJob, helpCrontabJobMu, robotNewJobsMap, helpActiveTimeHandler)
			CreateJob(ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap, defaultActiveTimeHandler)

			logMsg = "success"
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-dispatch").Info(logPrefix, logMsg, "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func JobActionProcess(jobType string, ch chan *Job, actionHandler func(*Job) (string, error), useCommon bool) {
	logPrefix := "JobActionProcess-" + jobType
	for {
		job := <-ch
		start := time.Now()
		logSufix := fmt.Sprintf("gid:%d,uid:%d", job.GuildID, job.UserID)
		logMsg := ""
		for {
			if actionHandler == nil {
				logMsg = "no action handler"
				break
			}
			//common process
			if useCommon {
				if err := robotActionBeforeCheck(job); err != nil {
					logMsg = "robotActionBeforeCheck Error: " + err.Error()
					break
				}
			}

			//specialized
			report, err := actionHandler(job)
			if err != nil {
				logMsg = "action failure: " + err.Error()
				break
			}
			logMsg = "action success: " + report
			break
		}
		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-action").Info(logPrefix, logMsg, logSufix, "cost:", cost, "ms")
	}
}

func defaultActiveTimeHandler() int64 {
	return time.Now().Unix() + int64(120+rand.Intn(301))
}
