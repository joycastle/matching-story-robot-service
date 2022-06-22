package club

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/joycastle/matching-story-robot-service/confmanager"
)

var (
	chatMsgMapping map[int][]string = make(map[int][]string)
)

func ReadRobotTeamChatFromConfManager() error {
	confMgr, err := confmanager.GetConfManagerVer().GetConfManager()
	if err != nil {
		return fmt.Errorf("confmanager initialization error:%s", err.Error())
	}

	//check robotTeamChat
	if num, err := confMgr.GetConfRobotTeamChatNum(); err != nil {
		return fmt.Errorf("confmanager read RobotTeamChat error:%s", err.Error())
	} else if num < 1 {
		return fmt.Errorf("confmanager RobotTeamChat no parameters")
	} else {
		for i := 0; i < num; i++ {
			msg, err := confMgr.GetConfRobotTeamChatByIndex(i)
			if err != nil {
				return fmt.Errorf("confmanager robotTeamChat contents error:%s", err.Error())
			}
			t := msg.GetType()
			s := msg.GetTextNum()
			if _, ok := chatMsgMapping[t]; !ok {
				chatMsgMapping[t] = make([]string, 0)
			}
			chatMsgMapping[t] = append(chatMsgMapping[t], s)
		}
	}

	_, ok1 := chatMsgMapping[1]
	_, ok2 := chatMsgMapping[2]
	if !ok1 || !ok2 {
		return fmt.Errorf("confmanager robotTeamChat contents must content type 1 and type 2")
	}

	return nil
}

func getChatMsgByRand(index int) string {
	msgs := chatMsgMapping[index]
	length := len(msgs)

	rand.Seed(time.Now().UnixNano())
	return msgs[rand.Intn(length)]
}
