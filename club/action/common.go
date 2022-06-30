package action

import (
	"errors"
	"fmt"

	"github.com/joycastle/matching-story-robot-service/service"
)

func robotActionBeforeCheck(job *Job) error {
	//判断是否被踢出工会
	if u, err := service.GetGuildInfoByIDAndUid(job.GuildID, job.UserID); err != nil {
		return err
	} else if u.GuildID <= 0 {
		return errors.New(fmt.Sprintf("already kick out the guild"))
	}
	return nil
}
