package library

import (
	"sync"
)

func DeleteJobs(targets map[string]*Job, mu *sync.Mutex, newTargets map[string]*Job) int {
	mu.Lock()
	defer mu.Unlock()

	//delete
	for k, _ := range targets {
		if _, ok := newTargets[k]; !ok {
			delete(targets, k)
		}
	}

	length := len(targets)

	return length
}

func CreateJobs(targets map[string]*Job, mu *sync.Mutex, newTargets map[string]*Job) int {
	mu.Lock()
	defer mu.Unlock()

	for k, job := range newTargets {
		if _, ok := targets[k]; !ok {
			targets[k] = job
		}
	}
	length := len(targets)

	return length
}
