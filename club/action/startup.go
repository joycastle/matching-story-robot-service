package action

import (
	"fmt"
	"sync"
	"time"

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
	JOB_TYPE_ACTIVITY     = "Activity"
	JOB_TYPE_DELETE_GUILD = "DeleteGuild"
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

	//delete guild
	deleteGuildChannel chan *library.Job = make(chan *library.Job, 100)

	//using for ownaction
	ownActionCrontabJob     map[string]*library.Job = make(map[string]*library.Job, capacityMap)
	ownActionCrontabJobMu   *sync.Mutex             = new(sync.Mutex)
	ownActionProcessChannel chan *library.Job       = make(chan *library.Job, capacityChannel)

	//data source for update robot action
	RobotActionUpdateCrontabJob   map[string]*library.Job = make(map[string]*library.Job, capacityMap)
	RobotActionUpdateCrontabJobMu *sync.Mutex             = new(sync.Mutex)
)

func DeleteJob(targets map[string]*library.Job, mu *sync.Mutex, newTargets map[string]*library.Job) int {
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

func JobKey(userID, guildID int64) string {
	return fmt.Sprintf("%d-%d", userID, guildID)
}

func Startup() {
	go UpdateRobotJobs(10)

	go library.TaskTimed(JOB_TYPE_FIRSTIN, firstInCrontabJob, firstInCrontabJobMu, firstInProcessChannel, nil, 10)
	go library.TaskTimed(JOB_TYPE_REQUEST, requestCrontabJob, requestCrontabJobMu, requestProcessChannel, nil, 10)
	go library.TaskTimed(JOB_TYPE_REQUEST_CHAT, requestChatCrontabJob, requestChatCrontabJobMu, requestChatProcessChannel, nil, 10)
	go library.TaskTimed(JOB_TYPE_HELP, helpCrontabJob, helpCrontabJobMu, helpProcessChannel, nil, 10)
	go library.TaskTimed(JOB_TYPE_OWN_AI, ownActionCrontabJob, ownActionCrontabJobMu, ownActionProcessChannel, cycleTimeHandlerOwnAi, 10)

	go library.TaskProcess(JOB_TYPE_FIRSTIN, firstInProcessChannel, firstInActionHandler)
	go library.TaskProcess(JOB_TYPE_REQUEST, requestProcessChannel, robotActionBeforeCheck, requestActionHandler)
	go library.TaskProcess(JOB_TYPE_REQUEST_CHAT, requestChatProcessChannel, robotActionBeforeCheck, requestChatActionHandler)
	go library.TaskProcess(JOB_TYPE_HELP, helpProcessChannel, robotActionBeforeCheck, helpActionHandler)
	go library.TaskProcess(JOB_TYPE_OWN_AI, ownActionProcessChannel, robotActionBeforeCheck, ownActionHandler)
	go library.TaskProcess(JOB_TYPE_DELETE_GUILD, deleteGuildChannel, deleteActionHandler)

	//更新配置周一0点
	go UpdateRobotConfigMonday(RobotActionUpdateCrontabJob, RobotActionUpdateCrontabJobMu)
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
				info.Failed().Step(42).Err(err)
				break
			}

			filterUserTypeMaps := make(map[int64]string, len(userTypes))
			for _, v := range userTypes {
				filterUserTypeMaps[v.UserID] = v.UserType
			}

			//过滤工会状态
			guildDeleteInfos, err := service.GetGuildInfosWithField(guildIDs, []string{"deleted_at"})
			if err != nil {
				info.Failed().Step(43).Err(err)
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

			robotNewJobsMap := make(map[string]*library.Job, len(robotUsers))
			for _, v := range robotUsers {
				index := JobKey(v.UserID, v.GuildID)
				robotNewJobsMap[index] = library.NewEmptyJob().SetGuildID(v.GuildID).SetUserID(v.UserID)
			}

			//delete job
			library.DeleteJobs(firstInCrontabJob, firstInCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(requestCrontabJob, requestCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(requestChatCrontabJob, requestChatCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(helpCrontabJob, helpCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap)
			library.DeleteJobs(RobotActionUpdateCrontabJob, RobotActionUpdateCrontabJobMu, robotNewJobsMap)

			//create job
			info.Set(JOB_TYPE_REQUEST, library.CreateJobs(requestCrontabJob, requestCrontabJobMu, robotNewJobsMap))
			info.Set(JOB_TYPE_HELP, library.CreateJobs(helpCrontabJob, helpCrontabJobMu, robotNewJobsMap))
			info.Set(JOB_TYPE_OWN_AI, library.CreateJobs(ownActionCrontabJob, ownActionCrontabJobMu, robotNewJobsMap))

			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-data").Info("UpdateRobotJobs", info.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func robotActionBeforeCheck(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	//判断是否被踢出工会
	if u, err := service.GetGuildInfoByIDAndUid(job.GuildID, job.UserID); err != nil {
		return info.Failed().Step(101).Err(err)
	} else if u.GuildID <= 0 {
		return info.Failed().Step(102).ErrString("guild_id must > 0")
	}
	return info.Success()
}
