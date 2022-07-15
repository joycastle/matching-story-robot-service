package create

import (
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

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

			okDataMap := make(map[int64]struct{}, len(list)/2)
			delDataMap := make(map[int64]struct{}, len(list)/2)

			for _, v := range list {
				if v.DeletedAt.Valid == true {
					delDataMap[v.ID] = struct{}{}
				} else {
					okDataMap[v.ID] = struct{}{}
				}
			}

			okDataLen := len(okDataMap)
			delDataLen := len(delDataMap)

			if okDataLen > 0 {
				info.Set("create_task_map_len", MergeJobs(createTaskCronMap, createTaskCronMu, okDataMap))
				info.Set("kick_task_map_len", MergeJobs(kickTaskCronMap, kickTaskCronMu, okDataMap))
			}

			info.Success().Set("total", len(list), "okDataLen", okDataLen, "delDataLen", delDataLen)
			break
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-create").Info("PullDatas", info.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func MergeJobs(targets map[int64]*Job, mu *sync.Mutex, newTargets map[int64]struct{}) int {
	mu.Lock()
	defer mu.Unlock()

	//delete
	for k, _ := range targets {
		if _, ok := newTargets[k]; !ok {
			delete(targets, k)
		}
	}
	//create
	for k, _ := range newTargets {
		if _, ok := targets[k]; !ok {
			targets[k] = &Job{GuildID: k, IsInit: true}
		}
	}

	length := len(targets)

	return length
}
