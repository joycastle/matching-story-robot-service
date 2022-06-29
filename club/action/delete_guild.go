package action

import (
	"fmt"

	"github.com/joycastle/matching-story-robot-service/service"
)

func deleteActionHandler(job *Job) (string, error) {
	if err := service.DeleteGuild(job.GuildID); err != nil {
		return fmt.Sprintf("DeleteFailed,RobotNum:%d", job.RobotNum), err
	}
	return fmt.Sprintf("DeleteOK,RobotNum:%d", job.RobotNum), nil
}
