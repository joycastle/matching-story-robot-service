package action

import (
	"errors"
	"time"

	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

func helpActiveTimeHandler() int64 {
	return time.Now().Unix() + config.GetStrengthHelpTimeByRand()
}

func helpActionHandler(job *Job) (string, error) {
	//未完成的帮助
	notCompleteRequest, err := service.GetGuildHelpRequestNotComplete(job.GuildID)
	if err != nil {
		return "", err
	}
	if len(notCompleteRequest) == 0 {
		return "", errors.New("no new request")
	}

	//过滤机器人用户
	users := []int64{}
	for _, v := range notCompleteRequest {
		if v.RequesterID > 0 {
			users = append(users, v.RequesterID)
		}
	}
	if len(users) == 0 {
		return "", errors.New("no new reasonable request")
	}
	userTypes, err := service.GetUserTypes(users)
	if err != nil {
		return "", err
	}
	normalUserMap := make(map[int64]struct{})
	for _, v := range userTypes {
		if v.UserType == service.USERTYPE_NORMAL {
			normalUserMap[v.UserID] = struct{}{}
		}
	}
	if len(normalUserMap) == 0 {
		return "", errors.New("no normal user request")
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
	respones, err := service.BatchGetResponesUsers(requestIds)
	if err != nil {
		return "", err
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
		return "", errors.New("has help all request")
	}

	//发送请求
	report := ""
	for _, v := range targets {
		if resp, err := service.SendRequestHelpRPC("HELP", job.UserID, job.GuildID, v.ID); err != nil {
			return "", err
		} else {
			report = report + " " + resp.Data
		}
	}

	return report, nil
}
