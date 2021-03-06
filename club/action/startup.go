package action

import (
	"fmt"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/flowcontrol"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

const (
	JOB_TYPE_FIRSTIN      = "FirstJoin"
	JOB_TYPE_REQUEST      = "Request"
	JOB_TYPE_REQUEST_CHAT = "RequestChat"
	JOB_TYPE_HELP         = "Help"
	JOB_TYPE_MONDAY       = "MondayUpdate"
	JOB_TYPE_OWN_AI       = "OwnAI"
)

var (
	capacityMap     int = 10000
	capacityChannel int = 2000

	//firstIn job other operation see action_first_in.go
	firstInCrontabJob     map[string]*library.Job = make(map[string]*library.Job, 1000)
	firstInCrontabJobMu   *sync.Mutex             = new(sync.Mutex)
	firstInProcessChannel chan *library.Job       = make(chan *library.Job, 1000)

	//request job other operation see action_request.go
	requestCrontabJob     map[string]*library.Job = make(map[string]*library.Job, capacityMap)
	requestCrontabJobMu   *sync.Mutex             = new(sync.Mutex)
	requestProcessChannel chan *library.Job       = make(chan *library.Job, capacityChannel)

	//request chat job other operation see action_request.go
	requestChatCrontabJob     map[string]*library.Job = make(map[string]*library.Job, 1000)
	requestChatCrontabJobMu   *sync.Mutex             = new(sync.Mutex)
	requestChatProcessChannel chan *library.Job       = make(chan *library.Job, 1000)

	//help job other operation see action_help.go
	helpCrontabJob     map[string]*library.Job = make(map[string]*library.Job, capacityMap)
	helpCrontabJobMu   *sync.Mutex             = new(sync.Mutex)
	helpProcessChannel chan *library.Job       = make(chan *library.Job, capacityChannel)

	//using for ownaction
	ownActionCrontabJob     map[string]*library.Job = make(map[string]*library.Job, capacityMap)
	ownActionCrontabJobMu   *sync.Mutex             = new(sync.Mutex)
	ownActionProcessChannel chan *library.Job       = make(chan *library.Job, capacityChannel)

	//data source for update robot action
	RobotActionUpdateCrontabJob     map[string]*library.Job = make(map[string]*library.Job, capacityMap)
	RobotActionUpdateCrontabJobMu   *sync.Mutex             = new(sync.Mutex)
	RobotActionUpdateProcessChannel chan *library.Job       = make(chan *library.Job, capacityChannel)
)

func JobKey(userID, guildID int64) string {
	return fmt.Sprintf("%d-%d", userID, guildID)
}

func JobKeyHandler(job *library.Job) string {
	return fmt.Sprintf("%d-%d", job.GetUserID(), job.GetGuildID())
}

func Startup() {
	go UpdateRobotJobs(10)

	go library.TaskTimed(JOB_TYPE_FIRSTIN, firstInCrontabJob, firstInCrontabJobMu, firstInProcessChannel, 10)
	go library.TaskTimed(JOB_TYPE_REQUEST, requestCrontabJob, requestCrontabJobMu, requestProcessChannel, 10)
	go library.TaskTimed(JOB_TYPE_REQUEST_CHAT, requestChatCrontabJob, requestChatCrontabJobMu, requestChatProcessChannel, 10)
	go library.TaskTimed(JOB_TYPE_HELP, helpCrontabJob, helpCrontabJobMu, helpProcessChannel, 10)
	go library.TaskTimed(JOB_TYPE_OWN_AI, ownActionCrontabJob, ownActionCrontabJobMu, ownActionProcessChannel, 10)
	go library.TaskTimed(JOB_TYPE_MONDAY, RobotActionUpdateCrontabJob, RobotActionUpdateCrontabJobMu, RobotActionUpdateProcessChannel, 60)

	go library.TaskProcess(JOB_TYPE_FIRSTIN, firstInProcessChannel, firstInActionHandler)
	go library.TaskProcess(JOB_TYPE_REQUEST, requestProcessChannel, robotActionBeforeCheck, requestActionHandler)
	go library.TaskProcess(JOB_TYPE_REQUEST_CHAT, requestChatProcessChannel, robotActionBeforeCheck, requestChatActionHandler)
	go library.TaskProcess(JOB_TYPE_HELP, helpProcessChannel, robotActionBeforeCheck, helpActionHandler)
	go library.TaskProcess(JOB_TYPE_OWN_AI, ownActionProcessChannel, robotActionBeforeCheck, ownActionHandler)
	go library.TaskProcess(JOB_TYPE_MONDAY, RobotActionUpdateProcessChannel, robotActionBeforeCheck, robotActionUpdateHandler)
}

func UpdateRobotJobs(t int) error {
	for {
		start := time.Now()
		info := lib.NewLogStructed()
		for {
			userMap, err := service.GetAllGuildUserMapInfos()
			if err != nil {
				info.Failed().Step(40).Err(err)
				break
			}

			if len(userMap) == 0 {
				info.Failed().Step(41).ErrString("userMap empty")
				break
			}

			//??????uid?????????????????????
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

			//????????????????????????
			userTypes, err := service.GetUserInfosWithField(uids, []string{"user_type"})
			if err != nil {
				info.Failed().Step(42).Err(err)
				break
			}

			filterUserTypeMaps := make(map[int64]string, len(userTypes))
			for _, v := range userTypes {
				filterUserTypeMaps[v.UserID] = v.UserType
			}

			//??????????????????
			guildDeleteInfos, err := service.GetGuildInfosWithField(guildIDs, []string{"deleted_at"})
			if err != nil {
				info.Failed().Step(43).Err(err)
				break
			}
			filterGuildDeleteMap := make(map[int64]bool, len(guildIDs))
			for _, v := range guildDeleteInfos {
				filterGuildDeleteMap[v.ID] = v.DeletedAt.Valid
			}

			//??????????????????????????????
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

			//????????????robot
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

			robotNewJobsMap := make(map[string]*library.Job, len(robotUsers))
			for _, v := range robotUsers {
				if _, hit := flowcontrol.IsHit("robot-service", fmt.Sprintf("%d", v.GuildID), v.GuildID); hit {
					index := JobKey(v.UserID, v.GuildID)
					robotNewJobsMap[index] = library.NewEmptyJob().SetGuildID(v.GuildID).SetUserID(v.UserID)
				}
			}

			//delete job
			library.DeleteJobs(JOB_TYPE_FIRSTIN, firstInCrontabJob, firstInCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(JOB_TYPE_REQUEST, requestCrontabJob, requestCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(JOB_TYPE_REQUEST_CHAT, requestChatCrontabJob, requestChatCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(JOB_TYPE_HELP, helpCrontabJob, helpCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(JOB_TYPE_OWN_AI, ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(JOB_TYPE_MONDAY, RobotActionUpdateCrontabJob, RobotActionUpdateCrontabJobMu, robotNewJobsMap)

			//create job
			library.CreateJobs(JOB_TYPE_REQUEST, requestCrontabJob, requestCrontabJobMu, robotNewJobsMap, requestActiveTimeHandler, JobKeyHandler)
			library.CreateJobs(JOB_TYPE_HELP, helpCrontabJob, helpCrontabJobMu, robotNewJobsMap, helpActiveTimeHandler, JobKeyHandler)
			library.CreateJobs(JOB_TYPE_OWN_AI, ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap, cycleTimeHandlerOwnAi, JobKeyHandler)
			library.CreateJobs(JOB_TYPE_MONDAY, RobotActionUpdateCrontabJob, RobotActionUpdateCrontabJobMu, robotNewJobsMap, robotActionUpdateTimeHandler, JobKeyHandler)

			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-data").Info("UpdateRobotJobs", info.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func robotActionBeforeCheck(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	//???????????????????????????
	if u, err := service.GetGuildInfoByIDAndUid(job.GuildID, job.UserID); err != nil {
		return info.Failed().Step(101).Err(err)
	} else if u.GuildID <= 0 {
		return info.Failed().Step(102).ErrString("guild_id must > 0")
	}
	return info.Success()
}
