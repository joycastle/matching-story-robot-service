package club

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/joycastle/matching-story-robot-service/confmanager"
	"github.com/joycastle/matching-story-robot-service/confmanager/csvauto"
	"github.com/joycastle/matching-story-robot-service/model"
)

var (
	configIRobotTeamConfig csvauto.IRobotTeamConfig
	configRobotNames       []string
)

//读取配置信息
func ReadRobotTeamConfigFromConfManager() error {
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
	//read robotTeamConfig -> generate_robot_num
	if configIRobotTeamConfig.GetGenerateRobotNumLen() != 2 {
		return fmt.Errorf("confmanager robotTeamConfig -> generate_robot_num only 2 parameter")
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

type GuildLevelInfo struct {
	AvgLevel int
	MaxLevel int
	MinLevel int
}

//获取工会关卡信息
func getLevelInfo(users []model.User) GuildLevelInfo {
	var (
		gli        GuildLevelInfo
		totalLevel int = 0
		max        int = 0
		min        int = 999999
	)

	if len(users) == 0 {
		return gli
	}

	for _, user := range users {
		totalLevel = totalLevel + user.UserLevel
		if min > user.UserLevel {
			min = user.UserLevel
		}
		if max < user.UserLevel {
			max = user.UserLevel
		}
	}

	avg := float64(totalLevel) / float64(len(users))

	gli.AvgLevel = int(math.Floor(avg))
	gli.MaxLevel = max
	gli.MinLevel = min

	return gli
}

//获取随机关卡数
func getLevelByRand(guildInfo model.Guild, userInfos []model.User) int {
	levelInfo := getLevelInfo(userInfos)

	cmin := configIRobotTeamConfig.GetLevelRangeByIndex(0)
	cmax := configIRobotTeamConfig.GetLevelRangeByIndex(1)
	step := cmax - cmin
	if step < 0 {
		step = step * -1
	}
	//随机值
	rand.Seed(time.Now().UnixNano())
	randStep := cmin + rand.Intn(step+1)

	//随机加减策略
	if rand.Intn(100) >= 50 {
		randStep = randStep * -1
	}

	newLevel := levelInfo.AvgLevel + randStep

	//边界判断
	if newLevel < int(guildInfo.LevelLimit) {
		newLevel = int(guildInfo.LevelLimit)
	} else if newLevel > levelInfo.MaxLevel {
		newLevel = levelInfo.MaxLevel
	}

	return newLevel
}

//获取随机检测时间
func getActiveTimeByRand() int64 {
	min := configIRobotTeamConfig.GetInitialTimeRangeByIndex(0)
	max := configIRobotTeamConfig.GetInitialTimeRangeByIndex(1)
	step := max - min
	if step < 0 {
		step = step * -1
	}
	rand.Seed(time.Now().UnixNano())
	return time.Now().Unix() + int64(min+rand.Intn(step))
}

//获取机器人上线A
func getRobotMaxLimitNum() int {
	return configIRobotTeamConfig.GetRobotNumMaxlimit()
}

//获取真是用户数量B
func getNormalUserNum() int {
	return configIRobotTeamConfig.GetTeammemberNumLimit()
}

//获取单词生成机器人的数量
func getGenerateRobotNumByRand() int {
	cmin := configIRobotTeamConfig.GetGenerateRobotNumByIndex(0)
	cmax := configIRobotTeamConfig.GetGenerateRobotNumByIndex(1)
	step := cmax - cmin
	if step < 0 {
		step = step * -1
	}
	//随机值
	rand.Seed(time.Now().UnixNano())
	return cmin + rand.Intn(step+1)
}
