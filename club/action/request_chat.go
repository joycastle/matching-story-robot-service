package action

import (
	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

func CreateRequestChatJob(userID, guildID int64) {
	job := &Job{
		GuildID:    guildID,
		UserID:     userID,
		ActionTime: faketime.Now().Unix() + int64(config.GetHelpTalkTimeGapByRand()),
	}

	k := JobKey(userID, guildID)

	requestChatCrontabJobMu.Lock()
	requestChatCrontabJob[k] = job
	requestChatCrontabJobMu.Unlock()
}

func requestChatActionHandler(job *Job) *Result {
	chatMsg := config.GetChatMsgByRand(2)
	respone, err := service.SendChatMessageRPC("requestChat", job.UserID, job.GuildID, chatMsg)
	if err != nil {
		return ErrorText(400).Detail(err.Error(), "chatMsg", chatMsg, "respone", respone)
	}
	return ActionSuccess().Detail(respone.Data, "chatMsg", chatMsg)
}
