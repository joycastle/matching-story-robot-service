package create

import (
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/lib"
)

func TaskProcess(businessType string, ch chan *Job, handler func(*Job) *lib.LogStructuredJson) {
	for {
		job := <-ch
		start := time.Now()
		info := lib.NewLogStructed().Set("Job", job)

		if handler == nil {
			info.Failed().Step(1).ErrString("handler is nil")
		} else {
			info.Merge(handler(job))
		}

		cost := faketime.Since(start).Nanoseconds() / 1000000
		log.Get("club-create").Info("TaskProcess", businessType, info.String(), "cost:", cost, "ms")
	}
}
