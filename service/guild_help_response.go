package service

import (
	"fmt"
	"strings"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
)

func GetGuildResponeInfosWithFiledsByHelpIDs(ids []int64, fileds []string) ([]model.GuildHelpResponse, error) {
	var out []model.GuildHelpResponse

	if len(ids) == 0 {
		return out, nil
	}

	listArraySlice := lib.ArraySliceInt64(ids, 50)
	newFileds := MergeFileds(fileds, "id", "help_id")

	sqlTpl := fmt.Sprintf("SELECT %s FROM `guild_help_response` WHERE `help_id` IN ? AND `id` > ? ORDER BY `id` ASC LIMIT ?;", strings.Join(newFileds, ","))

	for _, vs := range listArraySlice {
		if len(vs) == 0 {
			continue
		}
		minID := int64(0)
		limit := 500
		for {
			var ret []model.GuildHelpResponse
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
