package service

import (
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
	"gorm.io/gorm"
)

//获取未完成的帮助
func GetGuildHelpRequestNotComplete(guildID int64) ([]model.GuildHelpRequest, error) {
	var (
		list  []model.GuildHelpRequest
		limit int   = 1000
		minID int64 = 0
	)

	for {
		var rets []model.GuildHelpRequest
		if r := mysql.Get("default-slave").Where("guild_id = ? AND id > ? AND total != count", guildID, minID).Order("id ASC").Limit(limit).Find(&rets); r.Error != nil {
			return nil, r.Error
		} else {
			list = append(list, rets...)
		}

		if len(rets) >= limit {
			minID = rets[limit-1].ID
		} else {
			break
		}
	}

	return list, nil
}

//获取最大ID
func GetMaxGuildHelpRequestID() (int64, error) {
	var id int64
	var req model.GuildHelpRequest
	if r := mysql.Get("default-slave").Model(&req).Select("MAX(id)").Scan(&id); r.Error != nil {
		if r.Error != gorm.ErrRecordNotFound {
			return 0, nil
		} else {
			return -1, r.Error
		}
	}
	return id, nil
}

//获取某个时间之后数据的最小ID
func GetMinGuildHelpRequestIDByTimeAfter(t int64) (int64, error) {
	var reqs []model.GuildHelpRequest
	if r := mysql.Get("default-slave").Where("time >= ?", t).Order("id ASC").Limit(1).Find(&reqs); r.Error != nil {
		return -1, r.Error
	}
	if len(reqs) == 0 {
		return -1, gorm.ErrRecordNotFound
	}
	return reqs[0].ID, nil
}

func GetGuildHelpRequestInfoByID(id int64) (model.GuildHelpRequest, error) {
	var req model.GuildHelpRequest
	if r := mysql.Get("default-slave").First(&req, id); r.Error != nil {
		return req, r.Error
	}
	return req, nil
}

//获取id之后的所有数据(不包括id这条数据)
func BatchGetGuildHelpRequestInfoByAfterID(id int64) ([]model.GuildHelpRequest, error) {
	var (
		list  []model.GuildHelpRequest
		limit int   = 100
		minID int64 = id
	)

	for {
		var rets []model.GuildHelpRequest
		if r := mysql.Get("default-slave").Where("id > ?", minID).Limit(limit).Find(&rets); r.Error != nil {
			return nil, r.Error
		} else {
			list = append(list, rets...)
		}

		if len(rets) >= limit {
			minID = rets[limit-1].ID
		} else {
			break
		}
	}

	return list, nil
}

//获取id之后的所有数据总量(不包括id这条数据)
func BatchGetGuildHelpRequestCountByAfterID(id int64) (int, error) {
	var count int64
	var req model.GuildHelpRequest
	if r := mysql.Get("default-slave").Model(&req).Where("id > ?", id).Count(&count); r.Error != nil {
		return -1, r.Error
	}
	return int(count), nil
}

func UpdateGuildHelpRequestCountByID(id int64, num int16) error {
	var req model.GuildHelpRequest
	if r := mysql.Get("default-master").Model(&req).Where("id = ?", id).Limit(1).Update("count", num); r.Error != nil {
		return r.Error
	}
	return nil
}

//获取工会中最新发送请求时间
type GhrLatest struct {
	GuildID     int64
	RequesterID int64
	MaxTime     int64
}

func BatchGetLatestReqeustTimeByGuildIDs(guildIDs []int64) ([]GhrLatest, error) {
	var ret []GhrLatest

	sqlTpl := "SELECT `guild_id`, `requester_id`, max(`time`) AS MaxTime FROM `guild_help_request` WHERE `guild_id` in (?) AND `requester_id` > 0 GROUP BY `guild_id`, `requester_id`;"

	if r := mysql.Get("default-slave").Raw(sqlTpl, guildIDs).Scan(&ret); r.Error != nil {
		return nil, r.Error
	}

	return ret, nil
}
