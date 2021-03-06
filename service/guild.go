package service

import (
	"fmt"
	"strings"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
)

func DeleteGuild(id int64) error {
	var ret model.Guild
	if r := mysql.Get("default-master").Delete(&ret, id); r.Error != nil {
		return r.Error
	}
	return nil
}

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
		limit int   = 200
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

func GetGuildInfosWithField(ids []int64, fileds []string) ([]model.Guild, error) {
	var out []model.Guild

	if len(ids) == 0 {
		return out, nil
	}

	sliceSize := 1000
	var listArraySlice [][]int64
	listArraySlice = make([][]int64, len(ids)/sliceSize+1)

	for k, v := range ids {
		index := k / sliceSize
		listArraySlice[index] = append(listArraySlice[index], v)
	}

	filedMap := make(map[string]struct{})
	filedMap["id"] = struct{}{}
	for _, v := range fileds {
		filedMap[v] = struct{}{}
	}
	newFileds := []string{}
	for k, _ := range filedMap {
		newFileds = append(newFileds, k)
	}

	sqlTpl := fmt.Sprintf("SELECT %s FROM `guild` WHERE `id` IN ? LIMIT ?;", strings.Join(newFileds, ","))

	for _, vs := range listArraySlice {
		if len(vs) == 0 {
			continue
		}
		var ret []model.Guild
		if r := mysql.Get("default-slave").Raw(sqlTpl, vs, sliceSize).Scan(&ret); r.Error != nil {
			return out, r.Error
		}
		out = append(out, ret...)
	}

	return out, nil
}
