package action

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/casual-server-lib/util"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	weedDaysConfig map[time.Weekday]int = make(map[time.Weekday]int, 7)
)

func init() {
	weedDaysConfig[time.Sunday] = 7
	weedDaysConfig[time.Monday] = 1
	weedDaysConfig[time.Tuesday] = 2
	weedDaysConfig[time.Wednesday] = 3
	weedDaysConfig[time.Thursday] = 4
	weedDaysConfig[time.Friday] = 5
	weedDaysConfig[time.Saturday] = 6
}

func cycleTimeHandlerOwnAi(job *Job) (int64, error) {
	robotConfig, err := service.GetRobotForGuild(job.UserID)
	if err != nil {
		return 0, err
	}

	if robotConfig.ID <= 0 {
		return 0, errors.New(fmt.Sprintf("not found robot info"))
	}

	actionID := int(robotConfig.GroupID)

	if actionID <= 0 {
		return 0, errors.New(fmt.Sprintf("action_id not found from robot_team.csv"))
	}

	actionTimes := int(robotConfig.ActNum)
	ts, err := config.GetSleepTimeByActionTimesByRand(actionID, actionTimes)
	if err != nil {
		return 0, err
	}

	return time.Now().Unix() + int64(ts), nil
}

func ownActionHandler(job *Job) (string, error) {
	robotConfig, err := service.GetRobotForGuild(job.UserID)
	if err != nil {
		return "", err
	}

	if robotConfig.ID <= 0 {
		return "", errors.New(fmt.Sprintf("not found robot info"))
	}

	actionID := robotConfig.GroupID
	if actionID <= 0 {
		return "", errors.New(fmt.Sprintf("action_id not found from robot_team.csv"))
	}

	activeDaysMap, err := config.GetRobotActiveDaysByActionID(int(actionID))
	if err != nil {
		return "", err
	}

	if len(activeDaysMap) == 0 {
		return "", errors.New(fmt.Sprintf("active day config not found active_id:%d", actionID))
	}

	if _, ok := activeDaysMap[-1]; ok && len(activeDaysMap) == 1 {
		return "", errors.New("active day config is [-1], not take effect")
	}

	todayWeek := time.Now().Weekday()
	todayWeekInt := weedDaysConfig[todayWeek]
	if _, ok := activeDaysMap[todayWeekInt]; !ok {
		return "", errors.New(fmt.Sprintf("today:week:%d is not a active day, %v, active_id:%d", todayWeekInt, activeDaysMap, actionID))
	}

	//获取最高积分
	//1.获取工会成员
	uids, err := service.GetGuildUserIds(job.GuildID)
	if err != nil {
		return "", err
	}
	userInfos, err := service.GetUserInfosWithField(uids, []string{"user_level"})
	if err != nil {
		return "", err
	}
	normalUserMaxLevel := 0
	currentUserLevel := 0
	for _, v := range userInfos {
		if v.UserType == service.USERTYPE_NORMAL {
			if normalUserMaxLevel < v.UserLevel {
				normalUserMaxLevel = v.UserLevel
			}
		} else {
			if v.UserID == job.UserID {
				currentUserLevel = v.UserLevel
			}
		}
	}

	//rule1判断
	rule1Limit, err := config.GetRule1TargetByRand(int(actionID))
	if err != nil {
		return "", err
	}
	if (currentUserLevel - normalUserMaxLevel) >= rule1Limit {
		return "", fmt.Errorf("Rule1,robot userlevel :%d Exceed normal userlevelmax:%d, limit:%d,will sleep", currentUserLevel, normalUserMaxLevel, rule1Limit)
	}

	//rule2判断
	rule2Limit, err := config.GetRule2TargetByRand(int(actionID))
	if err != nil {
		return "", err
	}
	if rule2Limit <= 0 || currentUserLevel > rule2Limit {
		return "", fmt.Errorf("Rule2, robot userlevel :%d Exceed limit:%d,will sleep", currentUserLevel, rule2Limit)
	}

	//增加关卡
	step, err := config.GetStepByActionTimesByRand(int(actionID), int(robotConfig.ActNum))
	if err != nil {
		return "", err
	}

	if err := service.UpdateUserLevelByUid(job.UserID, step); err != nil {
		return "", err
	}

	//增加次数
	if err := service.UpdateRobotActiveNumByUid(job.UserID, 1); err != nil {
		return "", err
	}

	return "", nil
}

func UpdateRobotConfigMonday(targets map[string]*Job, mu *sync.Mutex) {
	for {
		now := time.Now()
		nowStamp := now.Unix()
		mondayStamp := util.WeekMondayTimestamp(time.Now())
		sunStamp := mondayStamp + 86400*7
		timeDuration := sunStamp - nowStamp
		time.Sleep(time.Duration(timeDuration) * time.Second)

		start := time.Now()
		temp := make(map[string]*Job, capacityMap)
		mu.Lock()
		for k, v := range targets {
			temp[k] = v
		}
		mu.Unlock()

		//用户信息
		uids := []int64{}
		for _, job := range temp {
			uids = append(uids, job.UserID)
		}
		userLevels, err := service.GetUserInfosWithField(uids, []string{"user_level"})
		if err != nil {
			continue
		}
		userLevelsMap := make(map[int64]int, len(uids))
		for _, u := range userLevels {
			userLevelsMap[u.UserID] = u.UserLevel
		}

		//机器人信息
		robotConfigs, err := service.GetRobotInfosWithField(uids, []string{"conf_id", "group_id"})
		if err != nil {
			continue
		}
		type RobotConf struct {
			UserType int32
			ActionID int64
		}
		robotConfigsMap := make(map[string]RobotConf, len(robotConfigs))
		for _, robot := range robotConfigs {
			robotConfigsMap[robot.RobotID] = RobotConf{UserType: robot.ConfID, ActionID: robot.GroupID}
		}

		//获取要更新的目标
		needUpdateConfigs := make(map[string]int64, len(robotConfigs))
		needUpdateUids := []int64{}
		for _, job := range temp {
			userLevel, ok := userLevelsMap[job.UserID]
			if !ok {
				continue
			}
			robotID := service.GetRobotIDByUid(job.UserID)
			robotConfig, ok := robotConfigsMap[robotID]
			if !ok {
				continue
			}
			actionID := config.GetRobotActionIDByRand(userLevel, robotConfig.UserType)
			if actionID == robotConfig.ActionID {
				continue
			}

			needUpdateConfigs[robotID] = actionID
			needUpdateUids = append(needUpdateUids, job.UserID)
		}

		//开始执行更新操作
		if len(needUpdateConfigs) > 0 {
			//更新活动次数为0
			if err := service.ResetRobotActNum(needUpdateUids); err != nil {
				continue
			}
			for robotID, actionID := range needUpdateConfigs {
				if err := service.UpdateRobotByRobotID(robotID, "group_id", actionID); err != nil {
					continue
				}
				time.Sleep(1000)
			}
		}

		cost := time.Since(start).Nanoseconds() / 1000000
		log.Get("club-dispatch").Info("RobotActionIDUpdate", "processNum:", len(needUpdateConfigs), "cost:", cost, "ms")
	}
}
