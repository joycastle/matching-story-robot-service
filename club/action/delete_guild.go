package action

import (
	"github.com/joycastle/matching-story-robot-service/service"
)

func deleteActionHandler(job *Job) *Result {
	if err := service.DeleteGuild(job.GuildID); err != nil {
		return ErrorText(104).Detail(err, "RobotNum", job.RobotNum)
	}
	return ActionSuccess().Detail("RobotNum", job.RobotNum)
}
