package create

import (
	"math/rand"
	"sync"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	deleteMap map[int64]struct{} = make(map[int64]struct{}, 1000)
	deleteMu  *sync.Mutex        = new(sync.Mutex)
)

func deleteGuildTimeHandler(job *library.Job) (int64, error) {
	u := faketime.Now().Unix()
	r := int64(rand.Intn(30))
	log.Get("club-timed").Debug("Delete", r, u, job.String())
	return u + r, nil
}

func deleteGuildLogicHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()

	deleteMu.Lock()
	if _, ok := deleteMap[job.GuildID]; ok {
		deleteMu.Unlock()
		return info.Failed().Step(34).ErrString("alread delete")
	}
	deleteMu.Unlock()

	//1.获取机器人和真实用户分布
	userDistributions, err := service.GetGuildUserTypeDistribution(job.GuildID)
	if err != nil {
		return info.Failed().Step(30).Err(err)
	}

	var robotUsers []model.User
	var normalUsers []model.User

	if v, ok := userDistributions[service.USERTYPE_CLUB_ROBOT_SERVICE]; ok {
		robotUsers = v
	}
	if v, ok := userDistributions[service.USERTYPE_NORMAL]; ok {
		normalUsers = v
	}

	if len(normalUsers) == 0 {
		if err := service.DeleteGuild(job.GuildID); err != nil {
			return info.Failed().Err(err).Step(32)
		}

		deleteMu.Lock()
		deleteMap[job.GuildID] = struct{}{}
		deleteMu.Unlock()

		return info.Success().Set("robotNum", len(robotUsers))
	}

	return info.Failed().Step(33).Set("normalNum", len(normalUsers), "robotNum", len(robotUsers))
}
