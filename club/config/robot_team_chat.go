package config

import (
	"math/rand"

	"github.com/joycastle/matching-story-robot-service/confmanager"
)

var (
	chatMsgMapping map[int][]string = make(map[int][]string)
)

func ReadRobotTeamChatFromConfManager() error {
	confMgr, err := confmanager.GetConfManagerVer().GetConfManager()
	if err != nil {
		return errConfManangerInit("RobotTeamChat", err)
	}

	//check robotTeamChat
	if num, err := confMgr.GetConfRobotTeamChatNum(); err != nil {
		return errConfManangerRead("RobotTeamChat", err)
	} else if num < 1 {
		return errLineNumEmpty("RobotTeamChat")
	} else {
		for i := 0; i < num; i++ {
			msg, err := confMgr.GetConfRobotTeamChatByIndex(i)
			if err != nil {
				return errConfManangerRead("RobotTeamChat", err)
			}
			t := msg.GetType()
			s := msg.GetTextNum()
			if _, ok := chatMsgMapping[t]; !ok {
				chatMsgMapping[t] = make([]string, 0)
			}
			chatMsgMapping[t] = append(chatMsgMapping[t], s)
		}
	}

	m1, ok1 := chatMsgMapping[1]
	m2, ok2 := chatMsgMapping[2]
	if !ok1 || !ok2 || len(m1) == 0 || len(m2) == 0 {
		return errDataResultEmpty("RobotTeamChat", "chatMsgMapping")
	}

	return nil
}

func GetChatMsgByRand(index int) string {
	msgs := chatMsgMapping[index]
	length := len(msgs)
	return msgs[rand.Intn(length)]
}
