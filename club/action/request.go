package action

import (
	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/util"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	freezingTime int64 = 3600 * 7 //冷冻时间
)

func requestActiveTimeHandler(job *library.Job) (int64, error) {
	return faketime.Now().Unix() + config.GetStrengthRequestByRand(), nil
}

func requestActionHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	//获取我的请求记录
	list, err := service.GetHelpRequestByGuildIDAndUserID(job.GuildID, job.UserID)
	if err != nil {
		return info.Failed().Step(70).Err(err)
	}

	timeFilter := util.TimeStamp("2022-06-01 00:00:00")
	requestRecord := []model.GuildHelpRequest{}
	for _, v := range list {
		if v.Time >= timeFilter {
			requestRecord = append(requestRecord, v)
		}
	}

	need := false
	if len(requestRecord) == 0 {
		need = true
	} else {
		var maxTime int64 = 0
		for _, v := range requestRecord {
			if maxTime < v.Time {
				maxTime = v.Time
			}
		}

		if faketime.Now().Unix()-maxTime > freezingTime {
			need = true
		}
	}

	if !need {
		return info.Failed().Step(71).Err(err).Set("freezingTime", "7hour")
	}

	if resp, err := service.SendRequestRPC("REQUEST", job.UserID, job.GuildID); err != nil {
		return info.Failed().Step(72).Err(err).Set("resp", resp)
	} else {
		//发送私信
		CreateRequestChatJob(job.UserID, job.GuildID)
		return info.Success().Set("resp", resp)
	}
}
