package club

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/joycastle/matching-story-robot-service/confmanager"
	"github.com/joycastle/matching-story-robot-service/confmanager/csvauto"
)

var (
	configIRobotTeamConfig csvauto.IRobotTeamConfig
	configRobotNames       []string
)

//读取配置信息
func ReadConfigFromConfManager() error {
	confMgr, err := confmanager.GetConfManagerVer().GetConfManager()
	if err != nil {
		return fmt.Errorf("confmanager initialization error:%s", err.Error())
	}

	//read robotTeamConfig
	if num, err := confMgr.GetConfRobotTeamConfigNum(); err != nil {
		return fmt.Errorf("confmanager read RobotTeamConfig error:%s", err.Error())
	} else if num > 1 {
		return fmt.Errorf("confmanager RobotTeamConfig only need 2 parameter")
	}

	if v, err := confMgr.GetConfRobotTeamConfigByIndex(0); err != nil {
		return fmt.Errorf("confmanager RobotTeamConfig example initialization error:%s", err.Error())
	} else {
		configIRobotTeamConfig = v
	}
	//read robotTeamConfig -> initial_time_range
	if configIRobotTeamConfig.GetInitialTimeRangeLen() != 2 {
		return fmt.Errorf("confmanager robotTeamConfig -> initial_time_range only 2 parameter")
	}
	//read robotTeamConfig -> initial_like
	if configIRobotTeamConfig.GetInitialLikeLen() != 2 {
		return fmt.Errorf("confmanager robotTeamConfig -> initial_like only 2 parameter")
	}
	// read RobotName
	if num, err := confMgr.GetConfRobotNameNum(); err != nil {
		return fmt.Errorf("confmanager read RobotName error:%s", err.Error())
	} else if num <= 0 {
		return fmt.Errorf("confmanager RobotName is empty")
	} else {
		//read RobotName to Array
		for i := 0; i < num; i++ {
			if oneRob, err := confMgr.GetConfRobotNameByIndex(i); err != nil {
				return fmt.Errorf("confmanager read robot name error:%s", err.Error())
			} else {
				configRobotNames = append(configRobotNames, oneRob.GetName())
			}
		}
	}

	return nil
}

//获取机器人名称，existsNames 排除掉的名字
func getRobotNameByRand(existsNames ...string) string {
	var rangelist []string
	var filterMap map[string]struct{}
	if len(existsNames) > 0 {
		filterMap = make(map[string]struct{})
		for _, name := range existsNames {
			filterMap[name] = struct{}{}
		}
	}
	for _, name := range configRobotNames {
		if filterMap != nil {
			if _, ok := filterMap[name]; ok {
				continue
			}
		}
		rangelist = append(rangelist, name)
	}

	if len(rangelist) <= 0 {
		return "Empty"
	}

	rand.Seed(time.Now().UnixNano())

	return rangelist[rand.Intn(len(rangelist))]
}

//获取机器人头像
func getRobotIconByRand() string {
	//头像索引客户端配置
	appMinIndex := 1
	appMaxIndex := 14
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", appMinIndex+rand.Intn(appMaxIndex))
}

//获取随机点赞数
func getLikeNumByRand() int {
	min := configIRobotTeamConfig.GetInitialLikeByIndex(0)
	max := configIRobotTeamConfig.GetInitialLikeByIndex(1)
	step := max - min
	if step < 0 {
		step = step * -1
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(step+1)
}
