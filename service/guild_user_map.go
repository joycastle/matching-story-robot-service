package service

import (
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
)

func GetGuildInfoByIDAndUid(id, uid int64) (model.GuildUserMap, error) {
	var info model.GuildUserMap
	if r := mysql.Get("default-slave").Where("user_id = ? AND guild_id = ?", uid, id).Limit(1).Find(&info); r.Error != nil {
		return info, r.Error
	}
	return info, nil
}

//获取club的用户uid
func GetGuildUserIds(guildID int64) ([]int64, error) {
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
func GetGuildUserTypeDistribution(guildID int64) (map[string][]model.User, error) {
	uids, err := GetGuildUserIds(guildID)
	if err != nil {
		return nil, err
	}

	users, err := GetUserInfosWithField(uids, []string{})
	if err != nil {
		return nil, err
	}

	m := make(map[string][]model.User)

	for _, u := range users {
		m[u.UserType] = append(m[u.UserType], u)
	}

	return m, nil
}

//获取所有的工会信息
func GetAllGuildUserMapInfos() ([]model.GuildUserMap, error) {
	var (
		list  []model.GuildUserMap
		limit int   = 1000
		minID int64 = 0
	)

	sqlTpl := "SELECT * FROM `guild_user_map` WHERE `user_id` > ? ORDER BY `user_id` ASC LIMIT ?;"

	for {
		var ret []model.GuildUserMap
		if r := mysql.Get("default-slave").Raw(sqlTpl, minID, limit).Scan(&ret); r.Error != nil {
			return nil, r.Error
		} else {
			list = append(list, ret...)
		}

		if len(ret) >= limit {
			minID = ret[limit-1].UserID
		} else {
			break
		}
	}

	return list, nil
}
