package action

import (
	"time"

	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

func CreateFirstInJob(userID, guildID int64) {
	job := &Job{
		GuildID: guildID,
		UserID:  userID,
		///ActionTime: time.Now().Unix() + int64(config.GetJoinTalkTimeGapByRand()),
		ActionTime: time.Now().Unix() + 10,
	}

	k := JobKey(userID, guildID)

	firstInCrontabJobMu.Lock()
	firstInCrontabJob[k] = job
	firstInCrontabJobMu.Unlock()
}

func firstInActionHandler(job *Job) (string, error) {
	chatMsg := config.GetChatMsgByRand(1)
	respone, err := service.SendChatMessageRPC("firstInGuildProcess", job.UserID, job.GuildID, chatMsg)
	if err != nil {
		return "", err
	}
	return respone.Data, nil
}
