package action

import (
	"time"

	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

func CreateRequestChatJob(userID, guildID int64) {
	job := &Job{
		GuildID: guildID,
		UserID:  userID,
		///ActionTime: time.Now().Unix() + int64(config.GetHelpTalkTimeGapByRand()),
		ActionTime: time.Now().Unix() + 10,
	}

	k := JobKey(userID, guildID)

	requestChatCrontabJobMu.Lock()
	requestChatCrontabJob[k] = job
	requestChatCrontabJobMu.Unlock()
}

func requestChatActionHandler(job *Job) (string, error) {
	chatMsg := config.GetChatMsgByRand(2)
	respone, err := service.SendChatMessageRPC("requestChat", job.UserID, job.GuildID, chatMsg)
	if err != nil {
		return "", err
	}
	return respone.Data, nil
}
