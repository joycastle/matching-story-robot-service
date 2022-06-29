package service

import (
	"fmt"
	"time"

	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/matching-story-robot-service/model"
)

func GetGuildResponeUidsByRequestID(id int64) ([]int64, error) {
	var (
		res  []model.GuildHelpResponse
		uids []int64
	)
	if r := mysql.Get("default-slave").Where("help_id = ?", id).Find(&res); r.Error != nil {
		return uids, r.Error
	}

	for _, v := range res {
		if v.ResponderID > 0 {
			uids = append(uids, v.ResponderID)
		}
	}
	return uids, nil
}

func GetGuildResoneRobotUserByRequestID(id int64) ([]model.User, error) {
	var robotUsers []model.User
	uids, err := GetGuildResponeUidsByRequestID(id)
	if err != nil {
		return robotUsers, err
	}

	m, err := GetUserInfoWithTypeByUids(uids)
	if err != nil {
		return robotUsers, err
	}

	if us, ok := m[USERTYPE_CLUB_ROBOT_SERVICE]; ok {
		return us, nil
	}

	return robotUsers, nil
}

func AddGuildHelpResone(requestID, rspUserID, reqUserID int64) (model.GuildHelpResponse, error) {
	var rsp model.GuildHelpResponse
	rsp.HelpID = requestID
	rsp.ResponderID = rspUserID
	rsp.RequesterID = reqUserID
	rsp.Time = time.Now().Unix()

	if err := mysql.Get("default-master").Create(&rsp); err.Error != nil || err.RowsAffected != 1 {
		return rsp, fmt.Errorf("%s affected:%d", err.Error, err.RowsAffected)
	}

	return rsp, nil
}

func BatchGetResponesUsers(requestIds []int64) ([]model.GuildHelpResponse, error) {
	var res []model.GuildHelpResponse

	sqlTpl := "SELECT help_id, responder_id FROM `guild_help_response` WHERE `help_id` IN ?;"

	if r := mysql.Get("default-slave").Raw(sqlTpl, requestIds).Scan(&res); r.Error != nil {
		return res, r.Error
	}

	return res, nil
}
