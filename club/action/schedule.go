package action

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

const (
	JOB_TYPE_FIRSTIN      = "FirstJoinGuild"
	JOB_TYPE_REQUEST      = "Request"
	JOB_TYPE_REQUEST_CHAT = "RequestChat"
	JOB_TYPE_HELP         = "Help"
	JOB_TYPE_ACTIVITY     = "Activity"
	JOB_TYPE_DELETE_GUILD = "DeleteGuild"
	JOB_TYPE_MONDAY       = "MondayUpdate"
	JOB_TYPE_OWN_AI       = "OwnAI"
	JOB_TYPE_LEAVE_GUILD  = "LeaveGuild"
)

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
	RobotActionUpdateJob   map[string]*Job = make(map[string]*Job, capacityMap)
	RobotActionUpdateJobMu *sync.Mutex     = new(sync.Mutex)

	//leave to guild for guild clear robot logic
	leaveGuildJob            map[string]*Job = make(map[string]*Job, 1000) //only for distributing requests to reduce server pressure
	leaveGuildJobMu          *sync.Mutex     = new(sync.Mutex)
	leaveGuildProcessChannel chan *Job       = make(chan *Job, capacityChannel)
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

func CreateJob(targets map[string]*Job, mu *sync.Mutex, newTargets map[string]*Job, actionTimeHandler func() (int64, *Result)) int {
	c := 0
	mu.Lock()
	for k, newJob := range newTargets {
		if _, ok := targets[k]; !ok {
			if actionTimeHandler == nil {
				actionTimeHandler = defaultActiveTimeHandler
			}
			newJob.ActionTime, _ = actionTimeHandler()
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

func CrontabGenerateJob(jobType string, targets map[string]*Job, mu *sync.Mutex, output chan *Job, cycleTimeHandler func(*Job) (int64, *Result), t int) {
	for {
		start := time.Now()
		now := faketime.Now().Unix()
		logger := NewCrontabLog(jobType)

		needProcess := []*Job{}
		needResetProcess := []*Job{}
		total := 0

		mu.Lock()
		for k, job := range targets {
			if now-job.ActionTime >= 0 {
				needProcess = append(needProcess, job)
				delete(targets, k)
			}
		}
		total = len(targets)
		mu.Unlock()

		needProcessNum := len(needProcess)
		if needProcessNum > 0 {
			//循环定时器
			if cycleTimeHandler != nil {
				//单独处理防止加锁时间过长
				for _, job := range needProcess {
					next, ar := cycleTimeHandler(job)
					if ar.Code != 0 {
						logger.AddExtra(ar)
						next, _ = defaultActiveTimeHandler()
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
					total = len(targets)
					mu.Unlock()
				}
			}

			for _, job := range needProcess {
				output <- job
			}
		}

		logger.SetTotal(total).SetNew(needProcessNum).SetReset(len(needResetProcess))

		cost := faketime.Since(start).Nanoseconds() / 1000000
		log.Get("club-dispatch").Info("Crontab", logger.String(), "cost:", cost, "ms")
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
	go CrontabGenerateJob(JOB_TYPE_LEAVE_GUILD, leaveGuildJob, leaveGuildJobMu, leaveGuildProcessChannel, nil, 10)

	go JobActionProcess(JOB_TYPE_FIRSTIN, firstInProcessChannel, firstInActionHandler, false)
	go JobActionProcess(JOB_TYPE_REQUEST, requestProcessChannel, requestActionHandler, true)
	go JobActionProcess(JOB_TYPE_REQUEST_CHAT, requestChatProcessChannel, requestChatActionHandler, true)
	go JobActionProcess(JOB_TYPE_HELP, helpProcessChannel, helpActionHandler, true)
	go JobActionProcess(JOB_TYPE_OWN_AI, ownActionProcessChannel, ownActionHandler, true)
	go JobActionProcess(JOB_TYPE_LEAVE_GUILD, leaveGuildProcessChannel, leaveGuildHandler, false)
	go JobActionProcess(JOB_TYPE_DELETE_GUILD, deleteGuildChannel, deleteActionHandler, false)

	//更新配置周一0点
	go UpdateRobotConfigMonday(RobotActionUpdateJob, RobotActionUpdateJobMu)
}

func UpdateRobotJobs(t int) error {
	for {
		start := time.Now()
		logger := NewUpdateJobLog()

		for {
			userMap, err := service.GetAllGuildUserMapInfos()
			if err != nil {
				logger.SetState(ErrorText(100).Detail("guild_user_name", err.Error()).String())
				break
			}

			if len(userMap) == 0 {
				logger.SetState(ErrorText(101).Detail("guild_user_name").String())
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
				logger.SetState(ErrorText(100).Detail("user_table", err.Error()).String())
				break
			}

			filterUserTypeMaps := make(map[int64]string, len(userTypes))
			for _, v := range userTypes {
				filterUserTypeMaps[v.UserID] = v.UserType
			}

			//过滤工会状态
			guildDeleteInfos, err := service.GetGuildInfosWithField(guildIDs, []string{"deleted_at"})
			if err != nil {
				logger.SetState(ErrorText(100).Detail("guild", err.Error()).String())
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
			logger.AddDeleteInfo(JOB_TYPE_FIRSTIN, DeleteJob(firstInCrontabJob, firstInCrontabJobMu, robotNewJobsMap))
			logger.AddDeleteInfo(JOB_TYPE_REQUEST, DeleteJob(requestCrontabJob, requestCrontabJobMu, robotNewJobsMap))
			logger.AddDeleteInfo(JOB_TYPE_REQUEST_CHAT, DeleteJob(requestChatCrontabJob, requestChatCrontabJobMu, robotNewJobsMap))
			logger.AddDeleteInfo(JOB_TYPE_HELP, DeleteJob(helpCrontabJob, helpCrontabJobMu, robotNewJobsMap))
			logger.AddDeleteInfo(JOB_TYPE_OWN_AI, DeleteJob(ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap))
			logger.AddDeleteInfo(JOB_TYPE_MONDAY, DeleteJob(RobotActionUpdateJob, RobotActionUpdateJobMu, robotNewJobsMap))
			logger.AddDeleteInfo(JOB_TYPE_LEAVE_GUILD, DeleteJob(leaveGuildJob, leaveGuildJobMu, robotNewJobsMap))

			//create job
			logger.AddCreateInfo(JOB_TYPE_REQUEST, CreateJob(requestCrontabJob, requestCrontabJobMu, robotNewJobsMap, requestActiveTimeHandler))
			logger.AddCreateInfo(JOB_TYPE_HELP, CreateJob(helpCrontabJob, helpCrontabJobMu, robotNewJobsMap, helpActiveTimeHandler))
			logger.AddCreateInfo(JOB_TYPE_OWN_AI, CreateJob(ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap, defaultActiveTimeHandler))

			//sync guild job
			UpdateGuildJob(robotNewJobsMap)
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-dispatch").Info("UpdateRobotJobs", logger.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func JobActionProcess(jobType string, ch chan *Job, actionHandler func(*Job) *Result, useCommon bool) {
	for {
		job := <-ch
		start := time.Now()
		logger := NewScheduleLog(job, jobType)
		for {
			if actionHandler == nil {
				logger.SetError("action handler not exits")
				break
			}

			var actionResult *Result

			//common process
			if useCommon {
				actionResult = robotActionBeforeCheck(job)
				if actionResult.Code != 0 {
					logger.SetResult(actionResult)
					break
				}
			}

			//specialized
			actionResult = actionHandler(job)
			logger.SetResult(actionResult)
			break
		}
		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-action").Info(logger.String(), "cost:", cost, "ms")
	}
}

func defaultActiveTimeHandler() (int64, *Result) {
	return faketime.Now().Unix() + int64(120+rand.Intn(301)), ActionSuccess()
}
