package config

import (
	"fmt"

	"github.com/joycastle/matching-story-robot-service/confmanager"
)

var (
	activeDayMap map[int]map[int]struct{} = make(map[int]map[int]struct{})
)

//读取配置信息
func ReadRobotTeamFromConfManager() error {
	confMgr, err := confmanager.GetConfManagerVer().GetConfManager()
	if err != nil {
		return fmt.Errorf("confmanager initialization error:%s", err.Error())
	}

	//read robotTeamInitial
	if num, err := confMgr.GetConfRobotTeamNum(); err != nil {
		return fmt.Errorf("confmanager read RobotTeam error:%s", err.Error())
	} else if num <= 0 {
		return fmt.Errorf("confmanager RobotTeam is empty")
	} else {
		//init data
		for i := 0; i < num; i++ {
			iface, err := confMgr.GetConfRobotTeamByIndex(i)
			if err != nil {
				return fmt.Errorf("confmanager RobotTeam initialization error:%s", err.Error())
			}

			id := iface.GetID()
			tmp := make(map[int]struct{})
			for j := 0; j < iface.GetActivityDayLen(); j++ {
				tmp[iface.GetActivityDayByIndex(j)] = struct{}{}
			}
			activeDayMap[id] = tmp
		}
	}

	return nil
}

func GetRobotActiveDaysByActionID(aid int) map[int]struct{} {
	if _, ok := activeDayMap[aid]; ok {
		return activeDayMap[aid]
	}
	return make(map[int]struct{})
}
