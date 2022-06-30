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
	robot.ActNum = 0

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

func UpdateRobotActiveNumByUid(uid int64, step int) error {

	robotID := fmt.Sprintf("club_%d", uid)

	if step == 0 {
		return fmt.Errorf("step is zero")
	}
	var robot model.Robot
	sql := fmt.Sprintf("UPDATE `robot_table` SET `act_num` = `act_num` + %d WHERE `robot_id` = ? LIMIT 1;", step)
	if r := mysql.Get("default-master").Debug().Raw(sql, robotID).Scan(&robot); r.Error != nil {
		return r.Error
	}
	return nil
}
