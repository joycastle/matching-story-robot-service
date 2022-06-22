package service

import (
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
)

func GetGuildInfoByGuildID(guildID int64) (model.Guild, error) {
	var ret model.Guild
	if r := mysql.Get("default-slave").First(&ret, guildID); r.Error != nil {
		return ret, r.Error
	}
	return ret, nil
}

func GetAllGuildInfoFromDB() ([]model.Guild, error) {
	var (
		list  []model.Guild
		limit int   = 100
		minID int64 = 0
	)

	sqlTpl := "SELECT * FROM `guild` WHERE `id` > ? ORDER BY `id` ASC LIMIT ?;"

	for {
		var ret []model.Guild
		if r := mysql.Get("default-slave").Raw(sqlTpl, minID, limit).Scan(&ret); r.Error != nil {
			return nil, r.Error
		} else {
			list = append(list, ret...)
		}

		if len(ret) >= limit {
			minID = ret[limit-1].ID
		} else {
			break
		}
	}

	return list, nil
}
