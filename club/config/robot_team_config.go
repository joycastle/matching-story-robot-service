package config

import (
	"fmt"
	"math"
	"math/rand"

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

	//read robotTeamConfig -> join_talk_timegap
	if configIRobotTeamConfig.GetJoinTalkTimegapLen() != 2 {
		return fmt.Errorf("confmanager robotTeamConfig -> join_talk_timegap only 2 parameter")
	}

	//read robotTeamConfig -> life_request_timegap
	if configIRobotTeamConfig.GetLifeRequestTimegapLen() != 2 {
		return fmt.Errorf("confmanager robotTeamConfig -> life_request_timegap only 2 parameter")
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
func GetRobotNameByRand(existsNames ...string) string {
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

	//rand.Seed(time.Now().UnixNano())

	return rangelist[rand.Intn(len(rangelist))]
}

//获取机器人头像
func GetRobotIconByRand() string {
	//头像索引客户端配置
	appMinIndex := 1
	appMaxIndex := 14
	//rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", appMinIndex+rand.Intn(appMaxIndex))
}

//获取随机点赞数
func GetLikeNumByRand() int {
	min := configIRobotTeamConfig.GetInitialLikeByIndex(0)
	max := configIRobotTeamConfig.GetInitialLikeByIndex(1)
	step := max - min
	if step < 0 {
		step = step * -1
	}
	//rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(step+1)
}

type GuildLevelInfo struct {
	AvgLevel int
	MaxLevel int
	MinLevel int
}

//获取工会关卡信息
func GetLevelInfo(users []model.User) GuildLevelInfo {
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
func GetLevelByRand(guildInfo model.Guild, userInfos []model.User) int {
	levelInfo := GetLevelInfo(userInfos)

	cmin := configIRobotTeamConfig.GetLevelRangeByIndex(0)
	cmax := configIRobotTeamConfig.GetLevelRangeByIndex(1)
	step := cmax - cmin
	if step < 0 {
		step = step * -1
	}
	//随机值
	//rand.Seed(time.Now().UnixNano())
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
func GetActiveTimeByRand() int64 {
	min := configIRobotTeamConfig.GetInitialTimeRangeByIndex(0)
	max := configIRobotTeamConfig.GetInitialTimeRangeByIndex(1)
	step := max - min
	if step < 0 {
		step = step * -1
	}
	//rand.Seed(time.Now().UnixNano())
	return int64(min + rand.Intn(step+1))
}

//获取机器人上线A
func GetRobotMaxLimitNum() int {
	return configIRobotTeamConfig.GetRobotNumMaxlimit()
}

//获取真是用户数量B
func GetNormalUserNum() int {
	return configIRobotTeamConfig.GetTeammemberNumLimit()
}

//获取单词生成机器人的数量
func GetGenerateRobotNumByRand() int {
	cmin := configIRobotTeamConfig.GetGenerateRobotNumByIndex(0)
	cmax := configIRobotTeamConfig.GetGenerateRobotNumByIndex(1)
	step := cmax - cmin
	if step < 0 {
		step = step * -1
	}
	//随机值
	//rand.Seed(time.Now().UnixNano())
	return cmin + rand.Intn(step+1)
}

//获取机器人帮助延迟时间
func GetStrengthHelpTimeByRand() int64 {
	cmin := configIRobotTeamConfig.GetHelpTimegapByIndex(0)
	cmax := configIRobotTeamConfig.GetHelpTimegapByIndex(1)
	step := cmax - cmin
	if step < 0 {
		step = step * -1
	}
	//随机值
	//rand.Seed(time.Now().UnixNano())
	return int64(cmin + rand.Intn(step+1))
}

//获取初次进入机器人问候语延迟
func GetJoinTalkTimeGapByRand() int {
	cmin := configIRobotTeamConfig.GetJoinTalkTimegapByIndex(0)
	cmax := configIRobotTeamConfig.GetJoinTalkTimegapByIndex(1)
	step := cmax - cmin
	if step < 0 {
		step = step * -1
	}
	//随机值
	//rand.Seed(time.Now().UnixNano())
	return cmin + rand.Intn(step+1)
}

//获取机器人请求体力延迟
func GetStrengthRequestByRand() int64 {
	cmin := configIRobotTeamConfig.GetLifeRequestTimegapByIndex(0)
	cmax := configIRobotTeamConfig.GetLifeRequestTimegapByIndex(1)
	step := cmax - cmin
	if step < 0 {
		step = step * -1
	}
	//随机值
	//rand.Seed(time.Now().UnixNano())
	return int64(cmin + rand.Intn(step+1))
}
