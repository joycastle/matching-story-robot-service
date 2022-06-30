package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
)

const (
	USERTYPE_NORMAL             = "normal"   //正常用户
	USERTYPE_CLUB_ROBOT_SERVICE = "robclubs" //工会机器人

	DEVICE_TYPE_CLUB_ROBOT = 9 //工会机器人设备类型标识

	COUNTRY_CN = "cn"
	COUNTRY_EN = "en"
)

//创建工会机器人
func CreateGuildRobotUser(name, headIcon string, likeCnt, level int) (model.User, error) {
	var u model.User
	u.UserName = name
	u.UserHeadIcon = headIcon
	u.UserLikeCount = uint(likeCnt)
	u.UserLevel = level

	u.UserCountryData = COUNTRY_CN
	u.DeviceType = DEVICE_TYPE_CLUB_ROBOT
	u.UserType = USERTYPE_CLUB_ROBOT_SERVICE

	u.UserHelp = 0
	u.CreateTime = time.Now().Unix()
	u.UpdateTime = u.CreateTime

	if id, err := lib.GenerateUserID(); err != nil {
		return u, err
	} else {
		u.UserID = id
	}

	u.AccountID = lib.Md5(fmt.Sprintf("%d", u.UserID))

	if err := mysql.Get("default-master").Create(&u); err.Error != nil || err.RowsAffected != 1 {
		return u, fmt.Errorf("%s affected:%d", err.Error, err.RowsAffected)
	}

	return u, nil
}

//获取用户信息
func GetUserInfoByUserID(uid int64) (model.User, error) {
	var user model.User
	if r := mysql.Get("default-slave").First(&user, uid); r.Error != nil {
		return user, r.Error
	}
	return user, nil
}

//获取用户信息-批量
func BatchGetUserInfoByUserID(uids []int64) ([]model.User, error) {
	var users []model.User
	if r := mysql.Get("default-slave").Where("user_id in (?)", uids).Find(&users); r.Error != nil {
		return users, r.Error
	}
	return users, nil
}

//获取用户分布
func GetUserInfoWithTypeByUids(uids []int64) (map[string][]model.User, error) {
	users, err := BatchGetUserInfoByUserID(uids)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]model.User)

	for _, u := range users {
		m[u.UserType] = append(m[u.UserType], u)
	}
	return m, nil
}

func GetUserTypes(uids []int64) ([]model.User, error) {
	var out []model.User

	if len(uids) == 0 {
		return out, nil
	}

	sliceSize := 1000
	var listArraySlice [][]int64
	listArraySlice = make([][]int64, len(uids)/sliceSize+1)

	for k, v := range uids {
		index := k / sliceSize
		listArraySlice[index] = append(listArraySlice[index], v)
	}

	sqlTpl := "SELECT user_id, user_type FROM `user_table` WHERE `user_id` IN ? LIMIT ?;"

	for _, vs := range listArraySlice {
		var ret []model.User
		if r := mysql.Get("default-slave").Raw(sqlTpl, vs, sliceSize).Scan(&ret); r.Error != nil {
			return out, r.Error
		}
		out = append(out, ret...)
	}

	return out, nil
}

func GetUserInfosWithField(uids []int64, fileds []string) ([]model.User, error) {
	var out []model.User

	if len(uids) == 0 {
		return out, nil
	}

	sliceSize := 1000
	var listArraySlice [][]int64
	listArraySlice = make([][]int64, len(uids)/sliceSize+1)

	for k, v := range uids {
		index := k / sliceSize
		listArraySlice[index] = append(listArraySlice[index], v)
	}

	filedMap := make(map[string]struct{})
	filedMap["user_id"] = struct{}{}
	filedMap["user_type"] = struct{}{}
	for _, v := range fileds {
		filedMap[v] = struct{}{}
	}
	newFileds := []string{}
	for k, _ := range filedMap {
		newFileds = append(newFileds, k)
	}

	sqlTpl := fmt.Sprintf("SELECT %s FROM `user_table` WHERE `user_id` IN ? LIMIT ?;", strings.Join(newFileds, ","))

	for _, vs := range listArraySlice {
		var ret []model.User
		if r := mysql.Get("default-slave").Debug().Raw(sqlTpl, vs, sliceSize).Scan(&ret); r.Error != nil {
			return out, r.Error
		}
		out = append(out, ret...)
	}

	return out, nil
}

func UpdateUserLevelByUid(uid int64, step int) error {
	if step == 0 {
		return fmt.Errorf("step is zero")
	}
	var u model.User
	sql := fmt.Sprintf("UPDATE `user_table` SET `user_level` = `user_level` + %d WHERE `user_id` = ? LIMIT 1;", step)
	if r := mysql.Get("default-master").Debug().Raw(sql, uid).Scan(&u); r.Error != nil {
		return r.Error
	}
	return nil
}
