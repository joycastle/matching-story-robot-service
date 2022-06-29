package action

import (
	"errors"
	"time"

	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	freezingTime int64 = 3600 * 7 //冷冻时间
)

func requestActiveTimeHandler() int64 {
	return time.Now().Unix() + config.GetStrengthRequestByRand()
}

func requestActionHandler(job *Job) (string, error) {
	//获取我的请求记录
	requestRecord, err := service.GetHelpRequestByGuildIDAndUserID(job.GuildID, job.UserID)
	if err != nil {
		return "", err
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

		if time.Now().Unix()-maxTime > freezingTime {
			need = true
		}
	}

	if !need {
		return "", errors.New("not exceeded freeze time")
	}

	if resp, err := service.SendRequestRPC("REQUEST", job.UserID, job.GuildID); err != nil {
		return "", err
	} else {
		//发送私信
		CreateRequestChatJob(job.UserID, job.GuildID)
		return resp.Data, nil
	}
}
