package library

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/lib"
)

func TaskTimed(businessType string, targets map[string]*Job, mu *sync.Mutex, ch chan *Job, timeHandler func(*Job) (int64, error), cycSec int) {
	for {
		start := time.Now()
		info := lib.NewLogStructed()

		var (
			total          int             = 0
			expiredJobs    map[string]*Job = make(map[string]*Job, 100)
			expiredJobsLen int             = 0
			proJobs        []*Job          = make([]*Job, 0)
			proJobsLen     int             = 0
			initJobsLen    int             = 0
			deleteJobsLen  int             = 0
		)

		//check expired jobs
		mu.Lock()
		for k, job := range targets {
			if job.Expired() {
				expiredJobs[k] = job
			}

			log.Get("club-timed").Debug(businessType, job.String(), job.GetActiveTimeDesc(), faketime.Now().Unix())
		}
		total = len(targets)
		mu.Unlock()

		expiredJobsLen = len(expiredJobs)

		//has expired task
		if expiredJobsLen > 0 {
			//cycle timed task
			if timeHandler != nil {
				//update atcion time from timehandler if return error using DefaultTimeHandler
				for _, job := range expiredJobs {
					actTime, err := timeHandler(job)
					if err != nil {
						info.Set(fmt.Sprintf("timeHandler:%v", job), err.Error())
						actTime, _ = DefaultTimeHandler(nil)
					}
					job.SetActiveTime(actTime)
				}
			}

			mu.Lock()
			for k, job := range expiredJobs {
				if _, ok := targets[k]; !ok {
					continue
				}
				if job.IsInit() {
					job.SetNormalStatus()
					initJobsLen++
				} else {
					proJobs = append(proJobs, job)
					if timeHandler == nil {
						delete(targets, k)
						deleteJobsLen++
					}
				}
			}
			total = len(targets)
			mu.Unlock()
		}

		proJobsLen = len(proJobs)

		if proJobsLen > 0 {
			for _, job := range proJobs {
				ch <- job
			}
		}

		info.Success().Set("total", total, "expired", expiredJobsLen, "pro", proJobsLen, "init", initJobsLen, "del", deleteJobsLen)

		cost := faketime.Since(start).Nanoseconds() / 1000000
		log.Get("club-timed").Info(businessType, info.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(cycSec) * time.Second)
	}
}

func DefaultTimeHandler(job *Job) (int64, error) {
	return faketime.Now().Unix() + int64(rand.Intn(180)+120), nil
}
