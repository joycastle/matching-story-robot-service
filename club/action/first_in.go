package action

import (
	"time"

	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

func CreateFirstInJob(userID, guildID int64) {
	job := &Job{
		GuildID:    guildID,
		UserID:     userID,
		ActionTime: time.Now().Unix() + int64(config.GetJoinTalkTimeGapByRand()),
	}

	k := JobKey(userID, guildID)

	firstInCrontabJobMu.Lock()
	firstInCrontabJob[k] = job
	firstInCrontabJobMu.Unlock()
}

func firstInActionHandler(job *Job) *Result {
	chatMsg := config.GetChatMsgByRand(1)
	respone, err := service.SendChatMessageRPC("firstInGuildProcess", job.UserID, job.GuildID, chatMsg)
	if err != nil {
		return ErrorText(400).Detail(err.Error(), "chatMsg", chatMsg, "respone", respone)
	}
	return ActionSuccess().Detail(respone.Data, "chatMsg", chatMsg)
}
