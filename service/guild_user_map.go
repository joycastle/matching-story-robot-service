package service

import (
	"fmt"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
)

//获取club的用户uid
func GetClubUserIds(guildID int64) ([]int64, error) {
	var users []model.GuildUserMap
	if r := mysql.Get("default-slave").Where("guild_id = ?", guildID).Find(&users); r.Error != nil {
		return nil, r.Error
	}

	var list []int64
	for _, v := range users {
		list = append(list, v.UserID)
	}

	return list, nil
}

//获取用户类型分布
func GetClubUserTypeDistribution(guildID int64) (map[string][]model.User, error) {
	uids, err := GetClubUserIds(guildID)
	if err != nil {
		return nil, err
	}

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

//加入工会
func JoinToClub(guildID, uid int64) error {
	var gum model.GuildUserMap
	gum.UserID = uid
	gum.GuildID = guildID
	if err := mysql.Get("default-master").Create(&gum); err.Error != nil || err.RowsAffected != 1 {
		return fmt.Errorf("%s affected:%d", err.Error, err.RowsAffected)
	}
	return nil
}
