package service

import (
	"fmt"
	"strings"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
)

func GetHelpRequestByGuildIDAndUserID(guildID, uid int64) ([]model.GuildHelpRequest, error) {
	var rets []model.GuildHelpRequest
	if r := mysql.Get("default-slave").Where("guild_id = ? AND requester_id = ?", guildID, uid).Find(&rets); r.Error != nil {
		return rets, r.Error
	}
	return rets, nil
}

func GetGuildRequestInfosWithFiledsByGuildIDs(ids []int64, fileds []string) ([]model.GuildHelpRequest, error) {
	var out []model.GuildHelpRequest

	if len(ids) == 0 {
		return out, nil
	}

	listArraySlice := lib.ArraySliceInt64(ids, 50)
	newFileds := MergeFileds(fileds, "id", "guild_id")

	sqlTpl := fmt.Sprintf("SELECT %s FROM `guild_help_request` WHERE `guild_id` IN ? AND `id` > ? ORDER BY `id` ASC LIMIT ?;", strings.Join(newFileds, ","))
	for _, vs := range listArraySlice {
		if len(vs) == 0 {
			continue
		}
		minID := int64(0)
		limit := 500
		for {
			var ret []model.GuildHelpRequest
			if r := mysql.Get("default-slave").Raw(sqlTpl, vs, minID, limit).Scan(&ret); r.Error != nil {
				return out, r.Error
			} else {
				out = append(out, ret...)
			}
			if len(ret) >= limit {
				minID = ret[limit-1].ID
			} else {
				break
			}
		}
	}

	return out, nil
}
