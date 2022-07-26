package config

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/joycastle/matching-story-robot-service/confmanager"
)

var (
	levelRangeMap map[string]int  = make(map[string]int)
	robotTypeMap  map[int][]int   = make(map[int][]int)
	robotWightMap map[int][]int   = make(map[int][]int)
	robotAiddMap  map[int][][]int = make(map[int][][]int)

	robotTypeIndexMap map[int]map[int]int = make(map[int]map[int]int)
)

//读取配置信息
func ReadRobotTeamInitialFromConfManager() error {
	confMgr, err := confmanager.GetConfManagerVer().GetConfManager()
	if err != nil {
		return errConfManangerInit("RobotTeamInitial", err)
	}

	//read robotTeamInitial
	if num, err := confMgr.GetConfRobotTeamInitialNum(); err != nil {
		return errConfManangerInit("RobotTeamInitial", err)
	} else if num <= 0 {
		return errLineNumEmpty("RobotTeamInitial")
	} else {
		//init data
		for i := 0; i < num; i++ {
			iface, err := confMgr.GetConfRobotTeamInitialByIndex(i)
			if err != nil {
				return errConfManangerRead("RobotTeamInitial", err)
			}

			length := iface.GetLevelRangeLen()
			if length <= 0 || length >= 3 {
				return errDataOnlyNeedLimit("RobotTeamInitial", i, "level_range", 1, 2)
			}

			v0 := 0
			v1 := 0

			if length == 1 {
				v0 = iface.GetLevelRangeByIndex(0)
				v1 = v0
			} else {
				v0 = iface.GetLevelRangeByIndex(0)
				v1 = iface.GetLevelRangeByIndex(1)
			}

			levelRangeMap[fmt.Sprintf("%d=%d", v0, v1)] = i

			for j := 0; j < iface.GetRobotTypeLen(); j++ {
				if _, ok := robotTypeMap[i]; !ok {
					robotTypeMap[i] = []int{}
				}
				robotTypeMap[i] = append(robotTypeMap[i], iface.GetRobotTypeByIndex(j))
			}

			for j := 0; j < iface.GetRobotWeightLen(); j++ {
				if _, ok := robotWightMap[i]; !ok {
					robotWightMap[i] = []int{}
				}
				robotWightMap[i] = append(robotWightMap[i], iface.GetRobotWeightByIndex(j))
			}

			for j := 0; j < iface.GetRobotAiIDLen(); j++ {
				if _, ok := robotAiddMap[i]; !ok {
					robotAiddMap[i] = [][]int{}
				}
				robotAiddMap[i] = append(robotAiddMap[i], iface.GetRobotAiIDByIndex(j))
			}
		}

		for k, vs := range robotTypeMap {
			if _, ok := robotTypeIndexMap[k]; !ok {
				robotTypeIndexMap[k] = make(map[int]int)
			}
			for index, v := range vs {
				robotTypeIndexMap[k][v] = index
			}
		}
	}

	m := make(map[int]struct{})
	m[len(levelRangeMap)] = struct{}{}
	m[len(robotTypeMap)] = struct{}{}
	m[len(robotWightMap)] = struct{}{}
	m[len(robotAiddMap)] = struct{}{}

	if len(m) != 1 {
		return errDataResultNotMatch("RobotTeamInitial", "level_range", "robot_type", "robot_weight", "robot_ai_d")
	}

	for i := 0; i < len(levelRangeMap); i++ {
		if len(robotTypeMap[i]) != len(robotWightMap[i]) || len(robotTypeMap[i]) != len(robotAiddMap[i]) {
			return errDataArrayNumNotMatch("RobotTeamInitial", i, "level_range", "robot_type", "robot_weight", "robot_ai_d")
		}
	}

	return nil
}

func getLevelRangeIndex(level int) int {
	for k, index := range levelRangeMap {
		arr := strings.Split(k, "=")
		min, _ := strconv.Atoi(arr[0])
		max, _ := strconv.Atoi(arr[1])
		if level > min && level <= max {
			return index
		}
	}
	return 0
}

func getWeightIndex(weights []int) int {
	m := make(map[string]int)
	lastMax := 0
	maxx := 0
	for index, v := range weights {
		min := lastMax
		max := min + v*10 - 1
		m[fmt.Sprintf("%d-%d", min, max)] = index
		lastMax = max + 1

		if maxx < lastMax {
			maxx = lastMax
		}
	}
	randNum := rand.Intn(maxx)

	for k, index := range m {
		arr := strings.Split(k, "-")
		min, _ := strconv.Atoi(arr[0])
		max, _ := strconv.Atoi(arr[1])
		if randNum >= min && randNum <= max {
			return index
		}
	}
	return 0
}

func getRobotTypeByIndex(index int) int {
	weights := robotWightMap[index]
	sindex := getWeightIndex(weights)
	return robotTypeMap[index][sindex]
}

func GetRobotActionIDByRand(level int, utype int32) int64 {
	index := getLevelRangeIndex(level)
	utypeIndex := robotTypeIndexMap[index][int(utype)]
	actionIds := robotAiddMap[index][utypeIndex]
	length := len(actionIds)
	return int64(actionIds[rand.Intn(length)])
}

func GetRobotTypeByRand(level int) int32 {
	index := getLevelRangeIndex(level)
	return int32(getRobotTypeByIndex(index))
}
