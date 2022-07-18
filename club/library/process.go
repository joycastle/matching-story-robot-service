package library

import (
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/lib"
)

func TaskProcess(businessType string, ch chan *Job, handlers ...func(*Job) *lib.LogStructuredJson) {
	for {
		job := <-ch
		start := time.Now()
		info := lib.NewLogStructed().Set("Job", job)

		if len(handlers) == 0 {
			info.Failed().Step(2).ErrString("handlers not set")
		} else {
			for _, handler := range handlers {
				if handler == nil {
					info.Failed().Step(3).ErrString("handler is nil")
					break
				}

				ret := handler(job)
				info.Merge(ret)
				if ret.IsFailed() {
					break
				}
			}
		}

		cost := faketime.Since(start).Nanoseconds() / 1000000
		log.Get("club-process").Info(businessType, info.String(), "cost:", cost, "ms")
	}
}
