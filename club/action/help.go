package action

import (
	"time"

	"github.com/joycastle/casual-server-lib/util"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

func helpActiveTimeHandler() (int64, *Result) {
	rnd := config.GetStrengthHelpTimeByRand()
	return time.Now().Unix() + rnd, ActionSuccess().Detail("rnd", rnd)
}

func helpActionHandler(job *Job) *Result {
	//未完成的帮助
	allRequest, err := service.GetGuildRequestInfosWithFiledsByGuildIDs([]int64{job.GuildID}, []string{"total", "count", "requester_id", "guild_id", "time"})
	if err != nil {
		return ErrorText(100).Detail("guild_help_request", err.Error())
	}

	timeFilter := util.TimeStamp("2022-06-01 00:00:00")
	var notCompleteRequest []model.GuildHelpRequest
	for _, v := range allRequest {
		if v.Total != v.Count && v.Time >= timeFilter {
			notCompleteRequest = append(notCompleteRequest, v)
		}
	}

	if len(notCompleteRequest) == 0 {
		return ErrorText(3000)
	}

	//过滤机器人用户
	users := []int64{}
	for _, v := range notCompleteRequest {
		if v.RequesterID > 0 {
			users = append(users, v.RequesterID)
		}
	}
	if len(users) == 0 {
		return ErrorText(3001)
	}
	userTypes, err := service.GetUserInfosWithField(users, []string{"user_type"})
	if err != nil {
		return ErrorText(100).Detail("user_table", err.Error(), users)
	}
	normalUserMap := make(map[int64]struct{})
	for _, v := range userTypes {
		if v.UserType == service.USERTYPE_NORMAL {
			normalUserMap[v.UserID] = struct{}{}
		}
	}
	if len(normalUserMap) == 0 {
		return ErrorText(3002)
	}

	//正常用户的请求
	normalRequests := []model.GuildHelpRequest{}
	requestIds := []int64{}
	for _, v := range notCompleteRequest {
		if _, ok := normalUserMap[v.RequesterID]; ok {
			normalRequests = append(normalRequests, v)
			requestIds = append(requestIds, v.ID)
		}
	}
	//查看请求帮助记录
	respones, err := service.GetGuildResponeInfosWithFiledsByHelpIDs(requestIds, []string{"responder_id"})
	if err != nil {
		return ErrorText(100).Detail("guild_help_respone", err.Error(), requestIds)
	}
	responesUserMap := make(map[int64]map[int64]struct{})
	for _, v := range respones {
		if _, ok := responesUserMap[v.HelpID]; !ok {
			responesUserMap[v.HelpID] = make(map[int64]struct{})
		}
		responesUserMap[v.HelpID][v.ResponderID] = struct{}{}
	}
	//获取要帮助的目标
	targets := []model.GuildHelpRequest{}
	for _, v := range normalRequests {
		if umap, ok := responesUserMap[v.ID]; ok {
			if _, okk := umap[job.UserID]; okk {
				continue
			}
		}
		targets = append(targets, v)
	}

	if len(targets) == 0 {
		return ErrorText(3003)
	}

	ret := ActionSuccess()
	//发送请求
	for _, v := range targets {
		if resp, err := service.SendRequestHelpRPC("HELP", job.UserID, job.GuildID, v.ID); err != nil {
			return ErrorText(400).Detail("help_id", v.ID)
		} else {
			ret.Detail("help_id", v.ID, resp.Data)
		}
	}

	return ret
}
