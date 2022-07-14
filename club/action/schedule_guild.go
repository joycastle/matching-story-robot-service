package action

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/log"
)

const (
	GUILD_JOB_TYPE_CLEAR_ROBOT_DISPATCH = "ClearRobotDispatch"
)

var (
	guildCapacityMap      int = 10000
	guildCcapacityChannel int = 2000

	clearRobotJobMap            map[int64]*GuildJob = make(map[int64]*GuildJob, guildCapacityMap)
	clearRobotJobMapMu          *sync.Mutex         = new(sync.Mutex)
	clearRobotJobProcessChannel chan *GuildJob      = make(chan *GuildJob, guildCcapacityChannel)
)

func GuildKeyFromJobKey(k string) int64 {
	arr := strings.Split(k, "-")
	v, _ := strconv.ParseInt(arr[1], 10, 64)
	return v
}

func getNewGuildTargetFromNewJobTargets(newTargets map[string]*Job) map[int64]*GuildJob {
	m := make(map[int64]*GuildJob, len(newTargets)/2)
	for k, _ := range newTargets {
		nk := GuildKeyFromJobKey(k)
		m[nk] = &GuildJob{GuildID: nk}
	}
	return m
}

func UpdateGuildJob(newTargets map[string]*Job) {
	newGuildTargets := getNewGuildTargetFromNewJobTargets(newTargets)

	//delete
	DeleteGuildJob(clearRobotJobMap, clearRobotJobMapMu, newGuildTargets)

	//create
	CreateGuildJob(clearRobotJobMap, clearRobotJobMapMu, newGuildTargets, nil)
}

func DeleteGuildJob(targets map[int64]*GuildJob, mu *sync.Mutex, newTargets map[int64]*GuildJob) int {
	c := 0
	mu.Lock()
	for k, _ := range targets {
		if _, ok := newTargets[k]; !ok {
			delete(targets, k)
			c = c + 1
		}
	}
	mu.Unlock()
	return c
}

func CreateGuildJob(targets map[int64]*GuildJob, mu *sync.Mutex, newTargets map[int64]*GuildJob, actionTimeHandler func() (int64, *Result)) int {
	c := 0
	mu.Lock()
	for k, newJob := range newTargets {
		if _, ok := targets[k]; !ok {
			if actionTimeHandler == nil {
				actionTimeHandler = defaultActiveTimeHandler
			}
			newJob.ActionTime, _ = actionTimeHandler()
			targets[k] = newJob
			c = c + 1
		}
	}
	mu.Unlock()
	return c
}

func CrontabGenerateGuildJob(jobType string, targets map[int64]*GuildJob, mu *sync.Mutex, output chan *GuildJob, cycleTimeHandler func(*GuildJob) (int64, *Result), t int) {
	for {
		start := time.Now()
		now := faketime.Now().Unix()
		logger := NewCrontabLog(jobType)

		needProcess := []*GuildJob{}
		needResetProcess := []*GuildJob{}
		total := 0

		mu.Lock()
		for k, job := range targets {
			if now-job.ActionTime >= 0 {
				needProcess = append(needProcess, job)
				delete(targets, k)
			}
		}
		total = len(targets)
		mu.Unlock()

		needProcessNum := len(needProcess)
		if needProcessNum > 0 {
			//循环定时器
			if cycleTimeHandler != nil {
				//单独处理防止加锁时间过长
				for _, job := range needProcess {
					next, ar := cycleTimeHandler(job)
					if ar.Code != 0 {
						logger.AddExtra(ar)
						next, _ = defaultActiveTimeHandler()
					}
					job.ActionTime = next
					needResetProcess = append(needResetProcess, job)
				}
				//重新设置时间
				if len(needResetProcess) > 0 {
					mu.Lock()
					for _, job := range needProcess {
						k := job.GuildID
						targets[k] = job
					}
					total = len(targets)
					mu.Unlock()
				}
			}

			for _, job := range needProcess {
				output <- job
			}
		}

		logger.SetTotal(total).SetNew(needProcessNum).SetReset(len(needResetProcess))

		cost := faketime.Since(start).Nanoseconds() / 1000000
		log.Get("club-guild-dispatch").Info("Crontab", logger.String(), "cost:", cost, "ms")
		time.Sleep(time.Duration(t) * time.Second)
	}
}

func StartupGuild() {
	go CrontabGenerateGuildJob(GUILD_JOB_TYPE_CLEAR_ROBOT_DISPATCH, clearRobotJobMap, clearRobotJobMapMu, clearRobotJobProcessChannel, nil, 10)
	go GuildJobActionProcess(GUILD_JOB_TYPE_CLEAR_ROBOT_DISPATCH, clearRobotJobProcessChannel, clearRobotDispatchHandler, false)
}

func GuildJobActionProcess(jobType string, ch chan *GuildJob, actionHandler func(*GuildJob) *Result, useCommon bool) {
	for {
		job := <-ch
		start := time.Now()
		logger := NewGuildScheduleLog(job, jobType)
		for {
			if actionHandler == nil {
				logger.SetError("action handler not exits")
				break
			}

			var actionResult *Result

			//common process
			if useCommon {
			}

			//specialized
			actionResult = actionHandler(job)
			logger.SetResult(actionResult)
			break
		}
		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-guild-action").Info(logger.String(), "cost:", cost, "ms")
	}
}
