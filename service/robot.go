package service

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
)

func GetRobotIDByUid(uid int64) string {
	return fmt.Sprintf("club_%d", uid)
}

func UpdateRobotRule2ConfigByUid(uid int64, target string) error {
	robotID := GetRobotIDByUid(uid)
	var robot model.Robot
	sql := "UPDATE `robot_table` SET `name` = ? WHERE `robot_id` = ? LIMIT 1;"
	if r := mysql.Get("default-master").Debug().Raw(sql, target, robotID).Scan(&robot); r.Error != nil {
		return r.Error
	}
	return nil
}

func CreateRobotForGuild(uid int64, utype int32, actionID int64, rule1SleepTime int, rule2SleepTime string) (model.Robot, error) {
	var robot model.Robot
	robot.RobotID = GetRobotIDByUid(uid)
	robot.ConfID = utype
	robot.GroupID = actionID
	robot.CreateTime = time.Now().Unix()
	robot.ActNum = 0
	robot.IndexNum = int32(rule1SleepTime)
	robot.HeadIcon = rule2SleepTime

	if err := mysql.Get("default-master").Create(&robot); err.Error != nil || err.RowsAffected != 1 {
		return robot, fmt.Errorf("%s affected:%d", err.Error, err.RowsAffected)
	}

	return robot, nil
}

func GetRobotForGuild(uid int64) (model.Robot, error) {
	robotID := GetRobotIDByUid(uid)

	var robot model.Robot
	if err := mysql.Get("default-slave").Where("robot_id = ?", robotID).Limit(1).Find(&robot); err.Error != nil {
		return robot, err.Error
	}

	return robot, nil
}

func UpdateRobotActiveNumByUid(uid int64, step int) error {

	robotID := GetRobotIDByUid(uid)

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

func GetRobotInfosWithField(uids []int64, fileds []string) ([]model.Robot, error) {
	var out []model.Robot

	if len(uids) == 0 {
		return out, nil
	}

	sliceSize := 1000
	var listArraySlice [][]string
	listArraySlice = make([][]string, len(uids)/sliceSize+1)

	for k, v := range uids {
		index := k / sliceSize
		listArraySlice[index] = append(listArraySlice[index], GetRobotIDByUid(v))
	}

	filedMap := make(map[string]struct{})
	filedMap["robot_id"] = struct{}{}
	for _, v := range fileds {
		filedMap[v] = struct{}{}
	}
	newFileds := []string{}
	for k, _ := range filedMap {
		newFileds = append(newFileds, k)
	}

	sqlTpl := fmt.Sprintf("SELECT %s FROM `robot_table` WHERE `robot_id` IN ? LIMIT ?;", strings.Join(newFileds, ","))

	for _, vs := range listArraySlice {
		var ret []model.Robot
		if r := mysql.Get("default-slave").Debug().Raw(sqlTpl, vs, sliceSize).Scan(&ret); r.Error != nil {
			return out, r.Error
		}
		out = append(out, ret...)
	}

	return out, nil
}

func ResetRobotWithFiled(uids []int64, args ...any) error {
	if len(uids) == 0 {
		return nil
	}

	sliceSize := 1000
	var listArraySlice [][]string
	listArraySlice = make([][]string, len(uids)/sliceSize+1)

	for k, v := range uids {
		index := k / sliceSize
		listArraySlice[index] = append(listArraySlice[index], GetRobotIDByUid(v))
	}

	setValues, err := MergeFiledsKV(args...)
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("UPDATE `robot_table` SET %s  WHERE `robot_id` IN ? LIMIT ?;", strings.Join(setValues, ","))

	for _, vs := range listArraySlice {
		var ret []model.Robot
		if r := mysql.Get("default-master").Debug().Raw(sql, vs, len(vs)).Scan(&ret); r.Error != nil {
			return r.Error
		}
	}

	return nil
}

func UpdateRobotByRobotID(robotID string, args ...interface{}) error {
	if len(args) == 0 || len(args)%2 != 0 {
		return fmt.Errorf("args error")
	}

	setValues := []string{}
	for i := 1; i < len(args); i += 2 {
		k := args[i-1]
		v := args[i]
		if reflect.TypeOf(k).Kind() != reflect.String {
			return fmt.Errorf("%s must be string", k)
		}
		if _, ok := TypesInt[reflect.TypeOf(v).Kind()]; !ok && reflect.TypeOf(v).Kind() != reflect.String {
			return fmt.Errorf("%s must be int int8 ....or string", v)
		}
		if reflect.TypeOf(v).Kind() == reflect.String {
			setValues = append(setValues, fmt.Sprintf("`%s`=%s", k, v))
		} else {
			setValues = append(setValues, fmt.Sprintf("`%s`=%d", k, v))
		}
	}

	sql := fmt.Sprintf("UPDATE `robot_table` SET %s  WHERE `robot_id` = ? LIMIT 1;", strings.Join(setValues, ","))

	var ret model.Robot
	if r := mysql.Get("default-master").Debug().Raw(sql, robotID).Scan(&ret); r.Error != nil {
		return r.Error
	}

	return nil
}
