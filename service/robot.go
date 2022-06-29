package service

import (
	"fmt"
	"time"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
)

func CreateRobotForGuild(uid int64, utype int32, actionID int64) (model.Robot, error) {
	var robot model.Robot
	robot.RobotID = fmt.Sprintf("club_%d", uid)
	robot.ConfID = utype
	robot.GroupID = actionID
	robot.CreateTime = time.Now().Unix()

	if err := mysql.Get("default-master").Create(&robot); err.Error != nil || err.RowsAffected != 1 {
		return robot, fmt.Errorf("%s affected:%d", err.Error, err.RowsAffected)
	}

	return robot, nil
}

func GetRobotForGuild(uid int64) (model.Robot, error) {
	robotID := fmt.Sprintf("club_%d", uid)

	var robot model.Robot
	if err := mysql.Get("default-slave").Where("robot_id = ?", robotID).Limit(1).Find(&robot); err.Error != nil {
		return robot, err.Error
	}

	return robot, nil
}
