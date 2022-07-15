package action

import (
	"math/rand"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/service"
)

func CreateLeaveGuildJob(userID, guildID int64) {
	job := &Job{
		GuildID:    guildID,
		UserID:     userID,
		ActionTime: faketime.Now().Unix() + int64(rand.Intn(20)),
	}

	k := JobKey(userID, guildID)

	leaveGuildJobMu.Lock()
	leaveGuildJob[k] = job
	leaveGuildJobMu.Unlock()
}

func leaveGuildHandler(job *Job) *Result {
	respone, err := service.SendLeaveGuildRPC("", job.UserID, job.GuildID)
	if err != nil {
		return ErrorText(400).Detail(err.Error(), "respone", respone)
	}
	return ActionSuccess().Detail(respone.Data)
}
