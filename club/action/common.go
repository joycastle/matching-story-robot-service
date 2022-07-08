package action

import (
	"github.com/joycastle/matching-story-robot-service/service"
)

func robotActionBeforeCheck(job *Job) *Result {
	//判断是否被踢出工会
	if u, err := service.GetGuildInfoByIDAndUid(job.GuildID, job.UserID); err != nil {
		return ErrorText(100).Detail("robotActionBeforeCheck", err)
	} else if u.GuildID <= 0 {
		return ErrorText(101).Detail("robotActionBeforeCheck")
	}
	return ActionSuccess().Detail("robotActionBeforeCheck")
}
