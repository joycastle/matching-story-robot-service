package action

import (
	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/util"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

func robotActionUpdateTimeHandler(job *library.Job) (int64, error) {
	now := faketime.Now()
	nowStamp := now.Unix()
	mondayStamp := util.WeekMondayTimestamp(now)
	sunStamp := mondayStamp + 86400*7
	timeDuration := sunStamp - nowStamp
	return faketime.Now().Unix() + timeDuration, nil
}

func robotActionUpdateHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()

	//用户信息
	users, err := service.GetUserInfosWithField([]int64{job.GetUserID()}, []string{"user_level"})
	if err != nil {
		return info.Failed().Step(711).Err(err)
	}
	if len(users) != 1 {
		return info.Failed().Step(712).ErrString("not found user")
	}

	user := users[0]

	//机器人信息
	robotConfigs, err := service.GetRobotInfosWithField([]int64{job.GetUserID()}, []string{"conf_id", "group_id"})
	if err != nil {
		return info.Failed().Step(713).Err(err)
	}
	if len(robotConfigs) != 1 {
		return info.Failed().Step(714).ErrString("robot user config not found")
	}

	robotConfig := robotConfigs[0]

	actionID := config.GetRobotActionIDByRand(user.UserLevel, robotConfig.ConfID)
	if actionID == robotConfig.GroupID {
		return info.Failed().Step(715).ErrString("same action id").Set("action_id", actionID)
	}

	rule1SleepTs, err := config.GetRule1TargetByRand(int(actionID))
	if err != nil {
		return info.Failed().Step(7151).Err(err)
	}

	rule2SleepTs, err := config.GetRule2TargetByRand(int(actionID))
	if err != nil {
		return info.Failed().Step(7152).Err(err)
	}

	if err := service.ResetRobotWithFiled([]int64{job.GetUserID()}, "act_num", 0, "name", "index_num", rule1SleepTs, "head_icon", rule2SleepTs); err != nil {
		return info.Failed().Step(716).Err(err)
	}

	robotID := service.GetRobotIDByUid(job.GetUserID())

	if err := service.UpdateRobotByRobotID(robotID, "group_id", actionID); err != nil {
		return info.Failed().Step(717).Err(err).Set("action_id", actionID, "robot_id", robotID)
	}

	return info.Success().Set("new_action_id", actionID, "old_action_id", robotConfig.GroupID, "robot_id", robotID)
}
