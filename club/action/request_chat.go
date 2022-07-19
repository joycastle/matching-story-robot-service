package action

import (
	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

func CreateRequestChatJob(userID, guildID int64) {
	k := JobKey(userID, guildID)
	requestChatCrontabJobMu.Lock()
	requestChatCrontabJob[k] = library.NewEmptyJob().SetGuildID(guildID).SetUserID(userID).SetActiveTime(faketime.Now().Unix() + int64(config.GetHelpTalkTimeGapByRand()))
	requestChatCrontabJobMu.Unlock()
}

func requestChatActionHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	chatMsg := config.GetChatMsgByRand(2)
	respone, err := service.SendChatMessageRPC("requestChat", job.UserID, job.GuildID, chatMsg)
	if err != nil {
		return info.Failed().Step(61).Err(err).Set("chatMsg", chatMsg, "respone", respone)
	}
	return info.Success().Set("chatMsg", chatMsg, "respone", respone)
}
