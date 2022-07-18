package create

import (
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

func addDeleteGuildTask(gid int64) {
	deleteGuildTaskChannel <- library.NewEmptyJob().SetGuildID(gid)
}

func deleteGuildLogicHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	if err := service.DeleteGuild(job.GuildID); err != nil {
		return info.Failed().Err(err).Step(30)
	}
	return info.Success()
}
