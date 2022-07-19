package action

import (
	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

func CreateFirstInJob(userID, guildID int64) {
	k := JobKey(userID, guildID)
	firstInCrontabJobMu.Lock()
	firstInCrontabJob[k] = library.NewEmptyJob().SetGuildID(guildID).SetUserID(userID).SetActiveTime(faketime.Now().Unix() + int64(config.GetJoinTalkTimeGapByRand()))
	firstInCrontabJobMu.Unlock()
}

func firstInActionHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	chatMsg := config.GetChatMsgByRand(1)
	respone, err := service.SendChatMessageRPC("firstInGuildProcess", job.UserID, job.GuildID, chatMsg)
	if err != nil {
		return info.Failed().Step(11).Err(err).Set("chatMsg", chatMsg, "respone", respone)
	}
	return info.Success().Set("chatMsg", chatMsg, "respone", respone)
}
