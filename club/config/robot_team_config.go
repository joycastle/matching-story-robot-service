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
		return errConfManangerInit("robotTeamConfig", err)
	}

	//read robotTeamConfig
	if num, err := confMgr.GetConfRobotTeamConfigNum(); err != nil {
		return errConfManangerRead("robotTeamConfig", err)
	} else if num != 1 {
		return errLineNumLimit("robotTeamConfig", 1)
	}

	if v, err := confMgr.GetConfRobotTeamConfigByIndex(0); err != nil {
		return errConfManangerRead("robotTeamConfig", err)
	} else {
		configIRobotTeamConfig = v
	}
	//read robotTeamConfig -> initial_time_range
	if configIRobotTeamConfig.GetInitialTimeRangeLen() != 2 {
		return errDataArrayNumLimit("robotTeamConfig", 1, "initial_time_range", 2)
	}
	//read robotTeamConfig -> initial_like
	if configIRobotTeamConfig.GetInitialLikeLen() != 2 {
		return errDataArrayNumLimit("robotTeamConfig", 1, "initial_like", 2)
	}
	//read robotTeamConfig -> generate_robot_num
	if configIRobotTeamConfig.GetGenerateRobotNumLen() != 2 {
		return errDataArrayNumLimit("robotTeamConfig", 1, "generate_robot_num", 2)
	}

	//read robotTeamConfig -> join_talk_timegap
	if configIRobotTeamConfig.GetJoinTalkTimegapLen() != 2 {
		return errDataArrayNumLimit("robotTeamConfig", 1, "join_talk_timegap", 2)
	}

	//read robotTeamConfig -> life_request_timegap
	if configIRobotTeamConfig.GetLifeRequestTimegapLen() != 2 {
		return errDataArrayNumLimit("robotTeamConfig", 1, "life_request_timegap", 2)
	}

	// read RobotName
	if num, err := confMgr.GetConfRobotNameNum(); err != nil {
		return errConfManangerRead("robotTeamConfig", err)
	} else if num <= 0 {
		return errLineNumEmpty("robotTeamConfig")
	} else {
		//read RobotName to Array
		for i := 0; i < num; i++ {
			if oneRob, err := confMgr.GetConfRobotNameByIndex(i); err != nil {
				return errConfManangerRead("robotTeamConfig", err)
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

	mina, maxa := Compare2Int(min, max)
	return mina + rand.Intn(maxa-mina+1)
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

	min := configIRobotTeamConfig.GetLevelRangeByIndex(0)
	max := configIRobotTeamConfig.GetLevelRangeByIndex(1)
	mina, maxa := Compare2Int(min, max)

	//随机值
	//rand.Seed(time.Now().UnixNano())
	randStep := mina + rand.Intn(maxa-mina+1)

	//随机加减策略
	if rand.Intn(100) >= 40 {
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

	mina, maxa := Compare2Int(min, max)
	return int64(mina + rand.Intn(maxa-mina+1))
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

	mina, maxa := Compare2Int(cmin, cmax)
	return mina + rand.Intn(maxa-mina+1)
}

//获取机器人帮助延迟时间
func GetStrengthHelpTimeByRand() int64 {
	cmin := configIRobotTeamConfig.GetHelpTimegapByIndex(0)
	cmax := configIRobotTeamConfig.GetHelpTimegapByIndex(1)

	mina, maxa := Compare2Int(cmin, cmax)
	return int64(mina + rand.Intn(maxa-mina+1))
}

//获取初次进入机器人问候语延迟
func GetJoinTalkTimeGapByRand() int {
	cmin := configIRobotTeamConfig.GetJoinTalkTimegapByIndex(0)
	cmax := configIRobotTeamConfig.GetJoinTalkTimegapByIndex(1)
	mina, maxa := Compare2Int(cmin, cmax)
	return mina + rand.Intn(maxa-mina+1)
}

//获取机器人请求体力延迟
func GetStrengthRequestByRand() int64 {
	cmin := configIRobotTeamConfig.GetLifeRequestTimegapByIndex(0)
	cmax := configIRobotTeamConfig.GetLifeRequestTimegapByIndex(1)

	mina, maxa := Compare2Int(cmin, cmax)
	return int64(mina + rand.Intn(maxa-mina+1))
}

//获取发送体力请求之后的延迟
func GetHelpTalkTimeGapByRand() int {
	cmin := configIRobotTeamConfig.GetHelpTalkTimegapByIndex(0)
	cmax := configIRobotTeamConfig.GetHelpTalkTimegapByIndex(1)
	mina, maxa := Compare2Int(cmin, cmax)
	return mina + rand.Intn(maxa-mina+1)
}
