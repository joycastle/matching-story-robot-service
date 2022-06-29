package action

import (
	"fmt"
	"sync"
	"time"

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
)

type Job struct {
	GuildID    int64
	UserID     int64
	ActionTime int64
}

var (
	//firstIn job other operation see action_first_in.go
	firstInCrontabJob     map[string]*Job = make(map[string]*Job, 1000)
	firstInCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	firstInProcessChannel chan *Job       = make(chan *Job, 1000)

	//request job other operation see action_request.go
	requestCrontabJob     map[string]*Job = make(map[string]*Job, 5000)
	requestCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	requestProcessChannel chan *Job       = make(chan *Job, 5000)

	//request chat job other operation see action_request.go
	requestChatCrontabJob     map[string]*Job = make(map[string]*Job, 5000)
	requestChatCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	requestChatProcessChannel chan *Job       = make(chan *Job, 5000)

	//help job other operation see action_help.go
	helpCrontabJob     map[string]*Job = make(map[string]*Job, 5000)
	helpCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	helpProcessChannel chan *Job       = make(chan *Job, 5000)

	//activity job other operation see action_activity.go
	activityCrontabJob     map[string]*Job = make(map[string]*Job, 5000)
	activityCrontabJobMu   *sync.Mutex     = new(sync.Mutex)
	activityProcessChannel chan *Job       = make(chan *Job, 5000)
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
				targets[k] = newJob
				c = c + 1
			}
		}
	}
	mu.Unlock()
	return c
}

func JobKey(userID, guildID int64) string {
	return fmt.Sprintf("%d-%d", userID, guildID)
}

func CrontabGenerateJob(jobType string, targets map[string]*Job, mu *sync.Mutex, output chan *Job, t int) {
	logPrefix := "CrontabGenerateJob-" + jobType
	for {
		start := time.Now()
		now := start.Unix()

		needProcess := []*Job{}

		mu.Lock()
		for k, job := range targets {
			if now-job.ActionTime >= 0 || k == "130714000009-125323777000079360" {
				needProcess = append(needProcess, job)
				delete(targets, k)
			}
		}
		count := len(targets)
		mu.Unlock()

		for _, job := range needProcess {
			output <- job
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-dispatch").Info(logPrefix, "total:", count, "need to process num:", len(needProcess), "cost:", cost, "ms")

		time.Sleep(time.Duration(t) * time.Second)
	}
}

func Startup() {
	go UpdateRobotJobs(10)

	go CrontabGenerateJob(JOB_TYPE_FIRSTIN, firstInCrontabJob, firstInCrontabJobMu, firstInProcessChannel, 10)
	go CrontabGenerateJob(JOB_TYPE_REQUEST, requestCrontabJob, requestCrontabJobMu, requestProcessChannel, 10)
	go CrontabGenerateJob(JOB_TYPE_REQUEST_CHAT, requestChatCrontabJob, requestChatCrontabJobMu, requestChatProcessChannel, 10)
	go CrontabGenerateJob(JOB_TYPE_HELP, helpCrontabJob, helpCrontabJobMu, helpProcessChannel, 10)
	go CrontabGenerateJob(JOB_TYPE_ACTIVITY, activityCrontabJob, activityCrontabJobMu, activityProcessChannel, 10)

	go JobActionProcess(JOB_TYPE_FIRSTIN, firstInProcessChannel, firstInActionHandler, false)
	//go JobActionProcess(JOB_TYPE_REQUEST, requestProcessChannel, nil)
	//go JobActionProcess(JOB_TYPE_REQUEST_CHAT, requestChatProcessChannel, nil)
	go JobActionProcess(JOB_TYPE_HELP, helpProcessChannel, helpActionHandler, true)
	//go JobActionProcess(JOB_TYPE_ACTIVITY, helpProcessChannel, nil)
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

			//获取uid和工会维度数据
			uids := []int64{}
			guildIDMaps := make(map[int64]struct{})
			for _, v := range userMap {
				uids = append(uids, v.UserID)
				guildIDMaps[v.GuildID] = struct{}{}
			}
			guildIDs := []int64{}
			for k, _ := range guildIDMaps {
				guildIDs = append(guildIDs, k)
			}

			//获取用户类型数据
			userTypes, err := service.GetUserTypes(uids)
			if err != nil {
				logMsg = "GetUserTypes Error: " + err.Error()
				break
			}

			filterUserTypeMaps := make(map[int64]string, len(userTypes))
			for _, v := range userTypes {
				filterUserTypeMaps[v.UserID] = v.UserType
			}

			//过滤工会状态
			guildDeleteInfos, err := service.GetGuildDeleteInfos(guildIDs)
			if err != nil {
				logMsg = "GetGuildDeleteInfos Error: " + err.Error()
				break
			}
			filterGuildDeleteMap := make(map[int64]bool, len(guildIDs))
			for _, v := range guildDeleteInfos {
				filterGuildDeleteMap[v.ID] = v.DeletedAt.Valid
			}

			//获取可用robot
			robotUsers := []model.GuildUserMap{}
			for _, v := range userMap {
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
			DeleteJob(activityCrontabJob, activityCrontabJobMu, robotNewJobsMap)

			//create job
			CreateJob(requestCrontabJob, requestCrontabJobMu, robotNewJobsMap, nil)
			CreateJob(helpCrontabJob, helpCrontabJobMu, robotNewJobsMap, helpActiveTimeHandler)
			CreateJob(activityCrontabJob, activityCrontabJobMu, robotNewJobsMap, nil)

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