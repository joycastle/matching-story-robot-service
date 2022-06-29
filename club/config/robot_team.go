package config

import (
	"fmt"
	"math/rand"

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
				return fmt.Errorf("confmanager RobotTeam data error ID:%d", id)
			}

			//activeStepMap
			tmpb := [][]int{}
			for j := 0; j < iface.GetLevelRangeLen(); j++ {
				tmpb = append(tmpb, iface.GetLevelRangeByIndex(j))
			}
			activeStepMap[id] = tmpb

			if len(tmpa) != len(tmpb) {
				return fmt.Errorf("confmanager RobotTeam data error leve_range ID:%d", id)
			}

			//activeSleepRule1Map
			tmpc := []int{}
			for j := 0; j < iface.GetRobotSleepRule1Len(); j++ {
				tmpc = append(tmpc, iface.GetRobotSleepRule1ByIndex(j))
			}
			if len(tmpc) != 2 {
				return fmt.Errorf("confmanager RobotTeam data error rule1 ID:%d", id)
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
				return fmt.Errorf("confmanager RobotTeam data error rule2 ID:%d", id)
			}
			activeSleepRule2TimeMap[id] = tmpe

		}
	}

	GetRule2TargetByRand(1009)

	return nil
}

//获取活跃天数
func GetRobotActiveDaysByActionID(aid int) map[int]struct{} {
	if _, ok := activeDayMap[aid]; ok {
		return activeDayMap[aid]
	}
	return make(map[int]struct{})
}

//根据次数获取延时时间
func GetSleepTimeByActionTimesByRand(aid, t int) int {
	arr := activeSleepTimeMap[aid][RangeIndexLORC(activeRangeTimesIndexMap[aid], t)]
	min, max := Compare2Int(arr[0], arr[1])
	return min + rand.Intn(max-min+1)
}

//根据次数获取步长
func GetStepByActionTimesByRand(aid, t int) int {
	arr := activeStepMap[aid][RangeIndexLORC(activeRangeTimesIndexMap[aid], t)]
	min, max := Compare2Int(arr[0], arr[1])
	return min + rand.Intn(max-min+1)
}

//获取rule1目标值
func GetRule1TargetByRand(aid int) int {
	arr := activeSleepRule1Map[aid]
	min, max := Compare2Int(arr[0], arr[1])
	return min + rand.Intn(max-min+1)
}

//获取rule2目标值
func GetRule2TargetByRand(aid int) int {
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
	total := (weekMax - weekMin) * 86400
	fmt.Println("aaaaaaaaa", weekMin, weekMax, total)
	return 0
}
