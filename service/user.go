package service

import (
	"fmt"
	"time"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
)

const (
	USERTYPE_NORMAL     = "normal"   //正常用户
	USERTYPE_CLUB_ROBOT = "robclubs" //工会机器人

	DEVICE_TYPE_CLUB_ROBOT = 9 //工会机器人设备类型标识

	COUNTRY_CN = "cn"
	COUNTRY_EN = "en"
)

//创建工会机器人
func CreateClubRobot(name, headIcon string, likeCnt, level int) (model.User, error) {
	var u model.User
	u.UserName = name
	u.UserHeadIcon = headIcon
	u.UserLikeCount = uint(likeCnt)
	u.UserLevel = level

	u.UserCountryData = COUNTRY_CN
	u.DeviceType = DEVICE_TYPE_CLUB_ROBOT
	u.UserType = USERTYPE_CLUB_ROBOT

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
