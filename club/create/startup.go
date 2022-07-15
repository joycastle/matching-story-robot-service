package create

import "sync"

const (
	JOB_TYPE_CREATE_ROBOT = "CreateRobot"
	JOB_TYPE_KICK_ROBOT   = "KickRobot"
)

var (
	//create robot
	createTaskChannel chan *Job      = make(chan *Job, 2000)
	createTaskCronMap map[int64]*Job = make(map[int64]*Job, 10000)
	createTaskCronMu  *sync.Mutex    = new(sync.Mutex)

	//create robot
	kickTaskChannel chan *Job      = make(chan *Job, 2000)
	kickTaskCronMap map[int64]*Job = make(map[int64]*Job, 10000)
	kickTaskCronMu  *sync.Mutex    = new(sync.Mutex)
)

func Startup() {
	go PullDatas(30)

	go TaskTimed(20, JOB_TYPE_CREATE_ROBOT, createTaskCronMap, createTaskCronMu, createTaskChannel, createRobotTimeHandler)
	go TaskTimed(20, JOB_TYPE_KICK_ROBOT, kickTaskCronMap, kickTaskCronMu, kickTaskChannel, kickRobotTimeHandler)

	go TaskProcess(JOB_TYPE_CREATE_ROBOT, createTaskChannel, createRobotLogicHandler)
	go TaskProcess(JOB_TYPE_KICK_ROBOT, kickTaskChannel, kickRobotLogicHandler)
}
