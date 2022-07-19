package create

import (
	"math/rand"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

func deleteGuildTimeHandler(*library.Job) (int64, error) {
	return faketime.Now().Unix() + int64(rand.Intn(30)), nil
}

func deleteGuildLogicHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()

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
	}

	return info.Success().Set("robotNum", len(robotUsers))
}
