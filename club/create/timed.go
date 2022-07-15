package create

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/lib"
)

func TaskTimed(t int, businessType string, targets map[int64]*Job, mu *sync.Mutex, ch chan *Job, timeHandler func(*Job) (int64, error)) {
	for {
		start := time.Now()
		now := faketime.Now().Unix()
		info := lib.NewLogStructed()

		total := 0
		initJobs := []*Job{}
		proJobs := []*Job{}

		mu.Lock()
		for k, job := range targets {
			if now-job.ActionTime > 0 {
				if job.IsInit {
					job.IsInit = false
					initJobs = append(initJobs, job)
					delete(targets, k)
				} else {
					proJobs = append(proJobs, job)
				}
			}
		}
		total = len(targets)
		mu.Unlock()

		initJobsLen := len(initJobs)
		proJobsLen := len(proJobs)

		if initJobsLen > 0 {
			if timeHandler == nil {
				timeHandler = defaultTimeHandler
			}

			for index, job := range initJobs {
				actTime, err := timeHandler(job)
				if err != nil {
					info.Set(fmt.Sprintf("%d", job.GuildID), err.Error())
					actTime, _ = defaultTimeHandler(nil)
				}
				initJobs[index].ActionTime = actTime
			}

			mu.Lock()
			for _, job := range initJobs {
				targets[job.GuildID] = job
			}
			total = len(targets)
			mu.Unlock()
		}

		if proJobsLen > 0 {
			for _, job := range proJobs {
				ch <- job
			}
		}

		info.Success().Set("total", total, "initJobsLen", initJobsLen, "proJobsLen", proJobsLen)

		cost := faketime.Since(start).Nanoseconds() / 1000000
		log.Get("club-create").Info("TaskTimed", businessType, info.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func defaultTimeHandler(job *Job) (int64, error) {
	return faketime.Now().Unix() + int64(rand.Intn(180)+120), nil
}
