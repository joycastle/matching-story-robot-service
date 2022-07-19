package library

import (
	"encoding/json"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
)

const (
	JobStatusInit   = 0
	JobStatusNormal = 1
)

type Job struct {
	ActionTime int64 `json:"action_time"`
	GuildID    int64 `json:"guild_id"`
	UserID     int64 `json:"user_id,omitempty"`
	Status     int8  `json:"status"`
}

func NewEmptyJob() *Job {
	return &Job{
		ActionTime: 0,
		Status:     JobStatusInit,
	}
}

func (job *Job) SetActiveTime(t int64) *Job {
	job.ActionTime = t
	return job
}

func (job *Job) GetActiveTime() int64 {
	return job.ActionTime
}

func (job *Job) GetActiveTimeDesc() string {
	return time.Unix(job.ActionTime, 0).Format("2006-01-02 15:04:05")
}

func (job *Job) SetGuildID(v int64) *Job {
	job.GuildID = v
	return job
}

func (job *Job) GetGuildID() int64 {
	return job.GuildID
}

func (job *Job) SetUserID(v int64) *Job {
	job.UserID = v
	return job
}

func (job *Job) GetUserID() int64 {
	return job.UserID
}

func (job *Job) SetNormalStatus() *Job {
	job.Status = JobStatusNormal
	return job
}

func (job *Job) IsInit() bool {
	if job.Status == JobStatusInit {
		return true
	}
	return false
}

func (job *Job) Expired() bool {
	if faketime.Now().Unix()-job.ActionTime >= 0 {
		return true
	}
	return false
}

func (job *Job) String() string {
	b, _ := json.Marshal(job)
	return string(b)
}
