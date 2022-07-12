package service

import (
	"fmt"
	"strings"

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
