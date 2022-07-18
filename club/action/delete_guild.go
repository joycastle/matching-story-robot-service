package action

import (
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

func deleteActionHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	if err := service.DeleteGuild(job.GuildID); err != nil {
		return info.Failed().Step(1).Err(err)
	}
	return info.Success()
}
