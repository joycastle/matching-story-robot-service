package action

import (
	"sync"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
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

func cycleTimeHandlerOwnAi(job *Job) (int64, *Result) {
	robotConfig, err := service.GetRobotForGuild(job.UserID)
	if err != nil {
		return 0, ErrorText(100).Detail("robot_table", err.Error())
	}

	if robotConfig.ID <= 0 {
		return 0, ErrorText(101).Detail("robot_table", job.UserID)
	}

	actionID := int(robotConfig.GroupID)

	if actionID <= 0 {
		return 0, ErrorText(103).Detail("group_id", "as action_id")
	}

	actionTimes := int(robotConfig.ActNum)
	ts, err := config.GetSleepTimeByActionTimesByRand(actionID, actionTimes)
	if err != nil {
		return 0, ErrorText(200).Detail("sleep time", err.Error(), "actionID", actionID, "actionTimes", actionTimes)
	}

	return faketime.Now().Unix() + int64(ts), ActionSuccess()
}

func ownActionHandler(job *Job) *Result {
	robotConfig, err := service.GetRobotForGuild(job.UserID)
	if err != nil {
		return ErrorText(100).Detail("robot_rable", err.Error())
	}

	if robotConfig.ID <= 0 {
		return ErrorText(101).Detail("robot_table", job.UserID)
	}

	actionID := robotConfig.GroupID
	if actionID <= 0 {
		return ErrorText(103).Detail("group_id", "as action_id")
	}

	activeDaysMap, err := config.GetRobotActiveDaysByActionID(int(actionID))
	if err != nil {
		return ErrorText(200).Detail("active_days", err.Error())
	}

	if len(activeDaysMap) == 0 {
		return ErrorText(201).Detail("active_days", actionID)
	}

	if _, ok := activeDaysMap[-1]; ok && len(activeDaysMap) == 1 {
		return ErrorText(1000).Detail("activeDaysMap", activeDaysMap)
	}

	todayWeek := faketime.Now().Weekday()
	todayWeekInt := weedDaysConfig[todayWeek]
	if _, ok := activeDaysMap[todayWeekInt]; !ok {
		return ErrorText(1001).Detail(
			"action_id", actionID,
			"todayWeek", todayWeek,
			"avtiveDays", activeDaysMap[todayWeekInt])
	}

	//获取最高积分
	//1.获取工会成员
	uids, err := service.GetGuildUserIds(job.GuildID)
	if err != nil {
		return ErrorText(100).Detail("guild_user_map", err.Error())
	}
	userInfos, err := service.GetUserInfosWithField(uids, []string{"user_level", "account_id"})
	if err != nil {
		return ErrorText(100).Detail("user_table", err.Error())
	}
	normalUserMaxLevel := 0
	currentUserLevel := 0
	userAccountID := ""
	for _, v := range userInfos {
		if v.UserType == service.USERTYPE_NORMAL {
			if normalUserMaxLevel < v.UserLevel {
				normalUserMaxLevel = v.UserLevel
			}
		} else {
			if v.UserID == job.UserID {
				currentUserLevel = v.UserLevel
				userAccountID = v.AccountID
			}
		}
	}

	if len(userAccountID) == 0 {
		return ErrorText(101).Detail("user_table", "account not found")
	}

	//rule1判断
	rule1Limit, err := config.GetRule1TargetByRand(int(actionID))
	if err != nil {
		return ErrorText(200).Detail("rule1", err.Error())
	}
	if (currentUserLevel - normalUserMaxLevel) >= rule1Limit {
		return ErrorText(1002).Detail(
			"currentUserLevel", currentUserLevel,
			"normalUserMaxLevel", normalUserMaxLevel,
			"rule1Limit", rule1Limit)
	}

	//rule2判断
	rule2Limit, err := config.GetRule2TargetByRand(int(actionID))
	if err != nil {
		return ErrorText(200).Detail("rule2", err.Error())
	}
	if rule2Limit <= 0 || currentUserLevel > rule2Limit {
		return ErrorText(1003).Detail(
			"currentUserLevel", currentUserLevel,
			"rule2Limit", rule2Limit)
	}

	//增加关卡
	step, err := config.GetStepByActionTimesByRand(int(actionID), int(robotConfig.ActNum))
	if err != nil {
		return ErrorText(200).Detail("LeveStep", err.Error(), "act_num", robotConfig.ActNum)
	}

	if err := service.UpdateUserLevelByUid(job.UserID, step); err != nil {
		return ErrorText(102).Detail("user_table", err.Error())
	}
	//增加次数
	if err := service.UpdateRobotActiveNumByUid(job.UserID, 1); err != nil {
		return ErrorText(102).Detail("robot_table", err.Error())
	}

	//增加积分
	rpcRet, err := service.SendUpdateScoreRPC(userAccountID, job.UserID, step)
	if err != nil {
		return ErrorText(400).Detail("update_score", err.Error())
	}

	return ActionSuccess().Detail("update_score_rpc", rpcRet)
}

func UpdateRobotConfigMonday(targets map[string]*Job, mu *sync.Mutex) {
	for {
		now := time.Now()
		nowStamp := now.Unix()
		mondayStamp := util.WeekMondayTimestamp(now)
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
