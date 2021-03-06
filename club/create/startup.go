package create

import (
	"fmt"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/flowcontrol"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

const (
	JOB_TYPE_CREATE_ROBOT = "CreateRobot"
	JOB_TYPE_KICK_ROBOT   = "KickRobot"
	JOB_TYPE_DELETE_GUILD = "DeleteGuild"
)

var (
	//create robot
	createTaskChannel chan *library.Job       = make(chan *library.Job, 2000)
	createTaskCronMap map[string]*library.Job = make(map[string]*library.Job, 10000)
	createTaskCronMu  *sync.Mutex             = new(sync.Mutex)

	//kick robot
	kickTaskChannel chan *library.Job       = make(chan *library.Job, 2000)
	kickTaskCronMap map[string]*library.Job = make(map[string]*library.Job, 10000)
	kickTaskCronMu  *sync.Mutex             = new(sync.Mutex)

	//delete guild
	deleteGuildTaskChannel chan *library.Job       = make(chan *library.Job, 10000)
	deleteGuildTaskCronMap map[string]*library.Job = make(map[string]*library.Job, 10000)
	deleteGuildTaskCronMu  *sync.Mutex             = new(sync.Mutex)
)

func Startup() {
	go PullDatas(30)

	go library.TaskTimed(JOB_TYPE_CREATE_ROBOT, createTaskCronMap, createTaskCronMu, createTaskChannel, 20)
	go library.TaskTimed(JOB_TYPE_KICK_ROBOT, kickTaskCronMap, kickTaskCronMu, kickTaskChannel, 20)
	go library.TaskTimed(JOB_TYPE_DELETE_GUILD, deleteGuildTaskCronMap, deleteGuildTaskCronMu, deleteGuildTaskChannel, 20)

	go library.TaskProcess(JOB_TYPE_CREATE_ROBOT, createTaskChannel, createRobotLogicHandler)
	go library.TaskProcess(JOB_TYPE_KICK_ROBOT, kickTaskChannel, kickRobotLogicHandler)
	go library.TaskProcess(JOB_TYPE_DELETE_GUILD, deleteGuildTaskChannel, deleteGuildLogicHandler)
}

func JobKey(id int64) string {
	return fmt.Sprintf("%d", id)
}
func JobKeyHandler(job *library.Job) string {
	return fmt.Sprintf("%d", job.GetGuildID())
}

func PullDatas(t int) {
	for {
		start := time.Now()
		info := lib.NewLogStructed()
		for {
			list, err := service.GetAllGuildInfoFromDB()
			if err != nil {
				info.Failed().Step(1).Err(err)
				break
			}

			okDataMap := make(map[string]*library.Job, len(list)/2)
			delDataLen := 0

			for _, v := range list {
				if v.DeletedAt.Valid == true {
					delDataLen++
				} else {
					if _, hit := flowcontrol.IsHit("robot-service", fmt.Sprintf("%d", v.ID), v.ID); hit {
						okDataMap[JobKey(v.ID)] = library.NewEmptyJob().SetGuildID(v.ID)
					}
				}
			}

			okDataLen := len(okDataMap)

			if okDataLen > 0 {
				library.DeleteJobs(JOB_TYPE_CREATE_ROBOT, createTaskCronMap, createTaskCronMu, okDataMap)
				library.DeleteJobs(JOB_TYPE_KICK_ROBOT, kickTaskCronMap, kickTaskCronMu, okDataMap)
				library.DeleteJobs(JOB_TYPE_DELETE_GUILD, deleteGuildTaskCronMap, deleteGuildTaskCronMu, okDataMap)

				library.CreateJobs(JOB_TYPE_CREATE_ROBOT, createTaskCronMap, createTaskCronMu, okDataMap, createRobotTimeHandler, JobKeyHandler)
				library.CreateJobs(JOB_TYPE_KICK_ROBOT, kickTaskCronMap, kickTaskCronMu, okDataMap, kickRobotTimeHandler, JobKeyHandler)
				library.CreateJobs(JOB_TYPE_DELETE_GUILD, deleteGuildTaskCronMap, deleteGuildTaskCronMu, okDataMap, deleteGuildTimeHandler, JobKeyHandler)
			}

			info.Success().Set("total", len(list), "new", okDataLen, "del", delDataLen)
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-data").Info("PullGuildDatas", info.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}
