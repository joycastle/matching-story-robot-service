package library

import (
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/lib"
)

func TaskTimed(businessType string, targets map[string]*Job, mu *sync.Mutex, ch chan *Job, cycSec int) {
	for {
		start := time.Now()
		info := lib.NewLogStructed()

		expiredJobs := []*Job{}

		mu.Lock()
		for k, job := range targets {
			if job.Expired() {
				expiredJobs = append(expiredJobs, job)
				delete(targets, k) //expired must delete
			}

			log.Get("club-timed").Debug(businessType, job.String(), job.GetActiveTimeDesc(), faketime.Now().Unix())
		}
		total := len(targets)
		mu.Unlock()

		for _, job := range expiredJobs {
			ch <- job
		}

		info.Success().Set("total", total, "expired", len(expiredJobs))

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-timed").Info(businessType, info.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(cycSec) * time.Second)
	}
}
