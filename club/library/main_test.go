package library

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestData(t *testing.T) {
	newTasks := make(map[string]*Job)
	for i := 10000; i < 19999; i++ {
		newTasks[fmt.Sprintf("%d", i)] = NewEmptyJob().SetGuildID(int64(i))
	}

	targets := make(map[string]*Job, 10000)
	mu := new(sync.Mutex)
	ch := make(chan *Job, 10000)

	if l := DeleteJobs(targets, mu, newTasks); l != 0 {
		t.Fatal("step=0", "DeleteJobs", l)
	}

	if l := CreateJobs(targets, mu, newTasks); l != 9999 {
		t.Fatal("step=1", "CreateJobs", l)
	}

	if l := DeleteJobs(targets, mu, newTasks); l != 9999 {
		t.Fatal("step=2", "DeleteJobs", l)
	}

	delete(newTasks, "19990")
	delete(newTasks, "19991")
	delete(newTasks, "19992")

	if l := CreateJobs(targets, mu, newTasks); l != 9999 {
		t.Fatal("step=4", "CreateJobs", l)
	}

	if l := DeleteJobs(targets, mu, newTasks); l != 9996 {
		t.Fatal("step=5", "DeleteJobs", l)
	}

	newTasks["20001"] = NewEmptyJob().SetGuildID(int64(20001))
	newTasks["20002"] = NewEmptyJob().SetGuildID(int64(20002))
	newTasks["20003"] = NewEmptyJob().SetGuildID(int64(20003))
	newTasks["20004"] = NewEmptyJob().SetGuildID(int64(20004))

	if l := DeleteJobs(targets, mu, newTasks); l != 9996 {
		t.Fatal("step=6", "DeleteJobs", l)
	}

	if l := CreateJobs(targets, mu, newTasks); l != 10000 {
		t.Fatal("step=7", "CreateJobs", l)
	}

	go TaskProcess("test", ch, nil)

	//TaskTimed("test", targets, mu, ch, nil, 1)
	TaskTimed("test", targets, mu, ch, func(job *Job) (int64, error) { return time.Now().Unix() + 10, nil }, 5)
}
