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

func GetUserInfosWithField(uids []int64, fileds []string) ([]model.User, error) {
	var out []model.User

	if len(uids) == 0 {
		return out, nil
	}

	sliceSize := 500
	listArraySlice := lib.ArraySliceInt64(uids, sliceSize)
	newFileds := MergeFileds(fileds, "id", "user_id", "user_type")

	sqlTpl := fmt.Sprintf("SELECT %s FROM `user_table` WHERE `user_id` IN ? LIMIT ?;", strings.Join(newFileds, ","))

	for _, vs := range listArraySlice {
		if len(vs) == 0 {
			continue
		}
		var ret []model.User
		if r := mysql.Get("default-slave").Raw(sqlTpl, vs, sliceSize).Scan(&ret); r.Error != nil {
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
	if r := mysql.Get("default-master").Raw(sql, uid).Scan(&u); r.Error != nil {
		return r.Error
	}
	return nil
}
