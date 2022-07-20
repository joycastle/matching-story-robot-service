package action

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
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

func cycleTimeHandlerOwnAi(job *library.Job) (int64, error) {
	robotConfig, err := service.GetRobotForGuild(job.UserID)
	if err != nil {
		return 0, err
	}

	if robotConfig.ID <= 0 {
		return 0, fmt.Errorf("confId not set")
	}

	actionID := int(robotConfig.GroupID)

	if actionID <= 0 {
		return 0, fmt.Errorf("actionId not set")
	}

	tts := 0
	actionTimes := int(robotConfig.ActNum)
	if actionTimes == 0 {
		ts, err := config.GetFirstActionTimeByRand(actionID)
		if err != nil {
			return 0, err
		}
		tts = ts
	} else {
		ts, err := config.GetSleepTimeByActionTimesByRand(actionID, actionTimes)
		if err != nil {
			return 0, err
		}
		tts = ts
	}

	return faketime.Now().Unix() + int64(tts), nil
}

func ownActionHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	robotConfig, err := service.GetRobotForGuild(job.UserID)
	if err != nil {
		return info.Failed().Step(51).Err(err)
	}

	if robotConfig.ID <= 0 {
		return info.Failed().Step(52).ErrString("confId is not config")
	}

	actionID := robotConfig.GroupID
	if actionID <= 0 {
		return info.Failed().Step(53).ErrString("actionID is not config")
	}

	activeDaysMap, err := config.GetRobotActiveDaysByActionID(int(actionID))
	if err != nil {
		return info.Failed().Step(54).Err(err)
	}

	if len(activeDaysMap) == 0 {
		return info.Failed().Step(55).ErrString("activeDaysMap empty")
	}

	if _, ok := activeDaysMap[-1]; ok && len(activeDaysMap) == 1 {
		return info.Failed().Step(56).ErrString("avtive not match").Set("activeDaysMap", activeDaysMap)
	}

	todayWeek := faketime.Now().Weekday()
	todayWeekInt := weedDaysConfig[todayWeek]
	if _, ok := activeDaysMap[todayWeekInt]; !ok {
		return info.Failed().Step(561).Set(
			"action_id", actionID,
			"todayWeek", todayWeek,
			"avtiveDays", activeDaysMap[todayWeekInt]).ErrString("not a active day")
	}

	//获取最高积分
	//1.获取工会成员
	uids, err := service.GetGuildUserIds(job.GuildID)
	if err != nil {
		return info.Failed().Step(57).Err(err)
	}
	userInfos, err := service.GetUserInfosWithField(uids, []string{"user_level", "account_id"})
	if err != nil {
		return info.Failed().Step(58).Err(err)
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
		return info.Failed().Step(59).ErrString("account not found")
	}

	//rule1判断
	rule1Limit := robotConfig.IndexNum
	if (currentUserLevel - normalUserMaxLevel) >= int(rule1Limit) {
		return info.Failed().Step(5900).Set(
			"currentUserLevel", currentUserLevel,
			"normalUserMaxLevel", normalUserMaxLevel,
			"rule1Limit", rule1Limit)
	}

	//rule2判断
	rule2Limit, index, err := config.GetRule2TargetByRand(int(actionID))
	if err != nil {
		return info.Failed().Step(591).Err(err)
	}

	staticUserLevel := 0
	if len(robotConfig.Name) == 0 {
		m := make(map[int]int)
		m[index] = currentUserLevel
		if err := UpdateRobotRule2Config(job.UserID, m); err != nil {
			return info.Failed().Step(5901).Err(err).Set("config index", index, "user_level", currentUserLevel)
		}
		staticUserLevel = currentUserLevel
	} else {
		m := make(map[int]int)
		if err := json.Unmarshal([]byte(robotConfig.Name), &m); err != nil {
			return info.Failed().Step(5902).Err(err).Set("config index", index, "user_level", currentUserLevel, "robotcfg", robotConfig.Name)
		}
		if staticLevel, ok := m[index]; !ok {
			m[index] = currentUserLevel
			if err := UpdateRobotRule2Config(job.UserID, m); err != nil {
				return info.Failed().Step(5902).Err(err).Set("config index", index, "user_level", currentUserLevel)
			}
			staticUserLevel = currentUserLevel
		} else {
			staticUserLevel = staticLevel
		}
	}

	if rule2Limit <= 0 || currentUserLevel > (rule2Limit+staticUserLevel) {
		return info.Failed().Step(5910).Set(
			"currentUserLevel", currentUserLevel,
			"rule2Limit", rule2Limit,
			"staticUserLevel", staticUserLevel)
	}

	//增加关卡
	step, err := config.GetStepByActionTimesByRand(int(actionID), int(robotConfig.ActNum))
	if err != nil {
		return info.Failed().Step(592).Err(err).Set("act_num", robotConfig.ActNum)
	}

	if err := service.UpdateUserLevelByUid(job.UserID, step); err != nil {
		return info.Failed().Step(593).Err(err)
	}
	//增加次数
	if err := service.UpdateRobotActiveNumByUid(job.UserID, 1); err != nil {
		return info.Failed().Step(594).Err(err)
	}

	//增加积分
	rpcRet, err := service.SendUpdateScoreRPC(userAccountID, job.UserID, step)
	if err != nil {
		return info.Failed().Step(595).Err(err).Set("resp", rpcRet)
	}

	return info.Success().Set("resp", rpcRet)
}

func UpdateRobotRule2Config(userId int64, m map[int]int) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := service.UpdateRobotRule2ConfigByUid(userId, string(b)); err != nil {
		return err
	}
	return nil
}
