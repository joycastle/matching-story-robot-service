package action

import (
	"errors"
	"fmt"
	"time"

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
	ts := config.GetSleepTimeByActionTimesByRand(actionID, actionTimes)

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

	activeDaysMap := config.GetRobotActiveDaysByActionID(int(actionID))

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
	rule1Limit := config.GetRule1TargetByRand(int(actionID))
	if (currentUserLevel - normalUserMaxLevel) >= rule1Limit {
		return "", fmt.Errorf("Rule1,robot userlevel :%d Exceed normal userlevelmax:%d, limit:%d,will sleep", currentUserLevel, normalUserMaxLevel, rule1Limit)
	}

	//rule2判断
	rule2Limit := config.GetRule2TargetByRand(int(actionID))
	if rule2Limit <= 0 || currentUserLevel > rule2Limit {
		return "", fmt.Errorf("Rule2, robot userlevel :%d Exceed limit:%d,will sleep", currentUserLevel, rule2Limit)
	}

	//增加关卡
	step := config.GetStepByActionTimesByRand(int(actionID), int(robotConfig.ActNum))

	if err := service.UpdateUserLevelByUid(job.UserID, step); err != nil {
		return "", err
	}

	//增加次数
	if err := service.UpdateRobotActiveNumByUid(job.UserID, 1); err != nil {
		return "", err
	}

	return "", nil
}
