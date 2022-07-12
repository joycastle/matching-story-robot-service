package config

import (
	"fmt"
	"math/rand"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/util"

	"github.com/joycastle/matching-story-robot-service/confmanager"
)

var (
	activeDayMap              map[int]map[int]struct{} = make(map[int]map[int]struct{})
	activeRangeTimesIndexMap  map[int]map[string]int   = make(map[int]map[string]int)
	activeSleepTimeMap        map[int][][]int          = make(map[int][][]int)
	activeStepMap             map[int][][]int          = make(map[int][][]int)
	activeSleepRule1Map       map[int][]int            = make(map[int][]int)
	activeSleepRule2TargetMap map[int][][]int          = make(map[int][][]int)
	activeSleepRule2TimeMap   map[int][]int            = make(map[int][]int)
)

//读取配置信息
func ReadRobotTeamFromConfManager() error {
	confMgr, err := confmanager.GetConfManagerVer().GetConfManager()
	if err != nil {
		return errConfManangerInit("RobotTeam", err)
	}

	//read robotTeamInitial
	if num, err := confMgr.GetConfRobotTeamNum(); err != nil {
		return errConfManangerRead("RobotTeam", err)
	} else if num <= 0 {
		return errLineNumEmpty("RobotTeam")
	} else {
		//init data
		for i := 0; i < num; i++ {
			iface, err := confMgr.GetConfRobotTeamByIndex(i)
			if err != nil {
				return errConfManangerRead("RobotTeam", err)
			}

			//activeDayMap init
			id := iface.GetID()
			tmp := make(map[int]struct{})
			for j := 0; j < iface.GetActivityDayLen(); j++ {
				tmp[iface.GetActivityDayByIndex(j)] = struct{}{}
			}
			activeDayMap[id] = tmp

			//activeRangeTimesIndexMap init
			tmpp := make(map[int]int)
			for j := 0; j < iface.GetActiviteRangeLen(); j++ {
				tmpp[iface.GetActiviteRangeByIndex(j)] = j
			}
			tmppp := ToRangeIndexLCR8(tmpp)
			activeRangeTimesIndexMap[id] = tmppp

			//activeSleepTimeMap init
			tmpa := [][]int{}
			for j := 0; j < iface.GetTimeGapLen(); j++ {
				tmpa = append(tmpa, iface.GetTimeGapByIndex(j))
			}
			activeSleepTimeMap[id] = tmpa

			if len(tmpa) != len(tmppp) {
				return errDataArrayNumNotMatch("RobotTeam", id, "time_gap", "avtive_range")
			}

			//activeStepMap
			tmpb := [][]int{}
			for j := 0; j < iface.GetLevelRangeLen(); j++ {
				tmpb = append(tmpb, iface.GetLevelRangeByIndex(j))
			}
			activeStepMap[id] = tmpb

			if len(tmpa) != len(tmpb) {
				return errDataArrayNumNotMatch("RobotTeam", id, "time_gap", "level_range")
			}

			//activeSleepRule1Map
			tmpc := []int{}
			for j := 0; j < iface.GetRobotSleepRule1Len(); j++ {
				tmpc = append(tmpc, iface.GetRobotSleepRule1ByIndex(j))
			}
			if len(tmpc) != 2 {
				return errDataArrayNumLimit("RobotTeam", id, "robot_sleep_rule1", 2)
			}
			activeSleepRule1Map[id] = tmpc

			//activeSleepRule2TargetMap
			tmpd := ParseStringType(iface.GetRobotSleepRule2Target())
			activeSleepRule2TargetMap[id] = tmpd

			//activeSleepRule2TimeMap
			tmpe := []int{}
			for j := 0; j < iface.GetRobotSleepRule2TimeLen(); j++ {
				tmpe = append(tmpe, iface.GetRobotSleepRule2TimeByIndex(j))
			}
			if len(tmpe) != len(tmpd) {
				return errDataArrayNumNotMatch("RobotTeam", id, "robot_sleep_rule2_time", "robot_sleep_rule2_target")
			}

			activeSleepRule2TimeMap[id] = tmpe

		}
	}

	return nil
}

//获取活跃天数
func GetRobotActiveDaysByActionID(aid int) (map[int]struct{}, error) {
	if _, ok := activeDayMap[aid]; ok {
		return activeDayMap[aid], nil
	}
	return nil, errParseIndexNotFound("robotTeam", "activeDayMap", fmt.Sprintf("%d", aid))
}

//根据次数获取延时时间
func GetSleepTimeByActionTimesByRand(aid, t int) (int, error) {
	arr, ok := activeSleepTimeMap[aid]
	if !ok {
		return 0, errParseIndexNotFound("robotTeam", "activeSleepTimeMap", fmt.Sprintf("%d", aid))
	}

	arr2, ok := activeRangeTimesIndexMap[aid]
	if !ok {
		return 0, errParseIndexNotFound("robotTeam", "activeRangeTimesIndexMap", fmt.Sprintf("%d", aid))
	}

	ret := arr[RangeIndexLORC(arr2, t)]
	min, max := Compare2Int(ret[0], ret[1])
	return min + rand.Intn(max-min+1), nil
}

//根据次数获取步长
func GetStepByActionTimesByRand(aid, t int) (int, error) {
	arr, ok := activeStepMap[aid]
	if !ok {
		return 0, errParseIndexNotFound("robotTeam", "activeStepMap", fmt.Sprintf("%d", aid))
	}

	arr2, ok := activeRangeTimesIndexMap[aid]
	if !ok {
		return 0, errParseIndexNotFound("robotTeam", "activeRangeTimesIndexMap", fmt.Sprintf("%d", aid))
	}

	k := RangeIndexLORC(arr2, t)
	ret := arr[k]

	//arr := activeStepMap[aid][RangeIndexLORC(activeRangeTimesIndexMap[aid], t)]
	min, max := Compare2Int(ret[0], ret[1])
	return min + rand.Intn(max-min+1), nil
}

//获取rule1目标值
func GetRule1TargetByRand(aid int) (int, error) {
	arr, ok := activeSleepRule1Map[aid]
	if !ok {
		return 0, errParseIndexNotFound("robotTeam", "activeSleepRule1Map", fmt.Sprintf("%d", aid))
	}
	min, max := Compare2Int(arr[0], arr[1])
	return min + rand.Intn(max-min+1), nil
}

//获取rule2目标值, 必须用GetRobotActiveDaysByActionID判断，才能准确
func GetRule2TargetByRand(aid int) (int, error) {
	v, ok := activeDayMap[aid]
	if !ok {
		return 0, errParseIndexNotFound("robotTeam", "activeDayMap", fmt.Sprintf("%d", aid))
	}
	if len(v) <= 1 {
		return 0, errParseIndexLimit("robotTeam", fmt.Sprintf("%d", aid), "activeDayMap", 2)
	}

	weekMin := 10
	weekMax := 0
	for week, _ := range activeDayMap[aid] {
		if weekMin > week {
			weekMin = week
		}
		if weekMax < week {
			weekMax = week
		}
	}

	now := faketime.Now()
	currentStamp := now.Unix()
	mondayStamp := util.WeekMondayTimestamp(now)
	disStamp := (weekMin - 1) * 86400
	filterFirstDay := mondayStamp + int64(disStamp)

	totalSeconds := (weekMax - weekMin + 1) * 86400
	proportionTime := RangeIndexWithSliceStep(activeSleepRule2TimeMap[aid], totalSeconds)
	index := -1
	//fmt.Println(proportionTime, weekMax, weekMin)
	//fmt.Println(util.FromUnixtime(currentStamp).Format("2006-01-02 15:04:05"))
	for k, v := range proportionTime {
		min, max := ValueWithRangeKey(k)
		//fmt.Println(util.FromUnixtime(filterFirstDay + int64(min)).Format("2006-01-02 15:04:05"))
		//fmt.Println(util.FromUnixtime(filterFirstDay + int64(max)).Format("2006-01-02 15:04:05"))
		if currentStamp >= filterFirstDay+int64(min) && currentStamp <= filterFirstDay+int64(max) {
			index = v
			break
		}
	}

	if index < 0 {
		return 0, errParseIndexNotFound("robotTeam", "activeSleepRule2TargetMap", fmt.Sprintf("%d", aid))
	}

	length := len(activeSleepRule2TargetMap[aid][index])

	return activeSleepRule2TargetMap[aid][index][rand.Intn(length)], nil
}
