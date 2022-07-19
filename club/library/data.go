package library

import (
	"math/rand"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/matching-story-robot-service/lib"
)

func DeleteJobs(businessType string, targets map[string]*Job, mu *sync.Mutex, newTargets map[string]*Job) {
	start := time.Now()
	info := lib.NewLogStructed()

	var (
		total  int = 0
		delNum int = 0
	)

	mu.Lock()
	//delete
	for k, _ := range targets {
		if _, ok := newTargets[k]; !ok {
			delete(targets, k)
			delNum++
		}
	}
	total = len(targets)
	mu.Unlock()

	info.Success().Set("total", total, "delete", delNum)

	cost := time.Since(start).Nanoseconds() / 1000000
	log.Get("club-merge").Info("DeleteJobs", businessType, info.String(), "cost:", cost, "ms")
}

func CreateJobs(businessType string, targets map[string]*Job, mu *sync.Mutex, newTargets map[string]*Job, timeHandler func(*Job) (int64, error), jobKeyHandler func(*Job) string) {
	start := time.Now()
	info := lib.NewLogStructed()

	initNewJobs := []*Job{}
	initJobs := []*Job{}
	total := 0

	mu.Lock()
	for k, job := range newTargets {
		if _, ok := targets[k]; !ok {
			initJobs = append(initJobs, job)
		}
	}
	total = len(targets)
	mu.Unlock()

	if len(initJobs) > 0 {
		for _, job := range initJobs {
			actTime, err := timeHandler(job)
			if err != nil {
				actTime, _ = DefaultTimeHandler(job)
				info.Set(job.String(), err.Error())
			}
			newJob := NewEmptyJob().SetGuildID(job.GetGuildID()).SetUserID(job.GetUserID()).SetActiveTime(actTime)
			initNewJobs = append(initNewJobs, newJob)
		}
	}

	if len(initNewJobs) > 0 {
		mu.Lock()
		for _, job := range initNewJobs {
			k := jobKeyHandler(job)
			targets[k] = job
		}
		total = len(targets)
		mu.Unlock()
	}

	info.Success().Set("total", total, "init", len(initNewJobs))

	cost := time.Since(start).Nanoseconds() / 1000000
	log.Get("club-merge").Info("CreateJobs", businessType, info.String(), "cost:", cost, "ms")
}

func DefaultTimeHandler(job *Job) (int64, error) {
	return faketime.Now().Unix() + int64(rand.Intn(180)+120), nil
}
