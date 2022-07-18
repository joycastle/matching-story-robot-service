package action

import (
	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/casual-server-lib/util"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

func helpActiveTimeHandler(job *library.Job) (int64, error) {
	return faketime.Now().Unix() + config.GetStrengthHelpTimeByRand(), nil
}

func helpActionHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	//未完成的帮助
	allRequest, err := service.GetGuildRequestInfosWithFiledsByGuildIDs([]int64{job.GuildID}, []string{"total", "count", "requester_id", "guild_id", "time"})
	if err != nil {
		return info.Failed().Step(31).Err(err)
	}

	timeFilter := util.TimeStamp("2022-06-01 00:00:00")
	var notCompleteRequest []model.GuildHelpRequest
	for _, v := range allRequest {
		if v.Total != v.Count && v.Time >= timeFilter {
			notCompleteRequest = append(notCompleteRequest, v)
		}
	}

	if len(notCompleteRequest) == 0 {
		return info.Failed().Step(32).ErrString("no complete request")
	}

	//过滤机器人用户
	users := []int64{}
	for _, v := range notCompleteRequest {
		if v.RequesterID > 0 {
			users = append(users, v.RequesterID)
		}
	}
	if len(users) == 0 {
		return info.Failed().Step(33).ErrString("no users")
	}
	userTypes, err := service.GetUserInfosWithField(users, []string{"user_type"})
	if err != nil {
		return info.Failed().Step(34).Err(err)
	}
	normalUserMap := make(map[int64]struct{})
	for _, v := range userTypes {
		if v.UserType == service.USERTYPE_NORMAL {
			normalUserMap[v.UserID] = struct{}{}
		}
	}
	if len(normalUserMap) == 0 {
		return info.Failed().Step(35).ErrString("no normal users")
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
		return info.Failed().Step(36).Err(err).Set("requestIds", requestIds)
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
		return info.Failed().Step(37).ErrString("no new targets")
	}

	//发送请求
	for _, v := range targets {
		if resp, err := service.SendRequestHelpRPC("HELP", job.UserID, job.GuildID, v.ID); err != nil {
			return info.Failed().Step(38).Err(err).Set("resp", resp, "helpID", v.ID)
		} else {
			info.Set("help_id", v.ID, resp.Data)
		}
	}

	return info.Success()
}
