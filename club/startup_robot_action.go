package club

import (
	"sync"
	"time"
)

const (
	//机器人行为定义
	ACTION_TYPE_FIRST            = "first"    //初次加入工会
	ACTION_TYPE_STRENGTH_REQUEST = "request"  //请求体力
	ACTION_TYPE_STRENGTH_HELP    = "help"     //帮助体力
	ACTION_TYPE_ACTIVITY         = "activity" //参与俱乐部活动
)

type Action struct {
	Type     string
	UserID   int64
	GuildID  int64
	JoinTime int64
}

var (
	actionCapacity = 4000

	//初次加入工会执行的操作
	actionChannel chan *Action      = make(chan *Action, actionCapacity)
	actionMaping  map[int64]*Action = make(map[int64]*Action, actionCapacity)
	actionMu      *sync.Mutex       = new(sync.Mutex)

	actionProcessNum int = 10 //行为处理并发数量
)

func StartupGuildRobotActions() {

	for i := 0; i < actionProcessNum; i++ {
		go actionProcess()
	}
}

func addFirstAction(guildID, uid int64) {
	actionChannel <- &Action{
		Type:     ACTION_TYPE_FIRST,
		UserID:   uid,
		GuildID:  guildID,
		JoinTime: time.Now().Unix(),
	}
}

func actionProcess() {
	for {
		action := <-actionChannel
		//usage like do{}while()
		for {
			switch action.Type {
			case ACTION_TYPE_FIRST:
				firstProcess(action)
			case ACTION_TYPE_STRENGTH_REQUEST:
			case ACTION_TYPE_STRENGTH_HELP:
			case ACTION_TYPE_ACTIVITY:
			}
		}
	}
}

func firstProcess(action *Action) error {
	return nil
}
