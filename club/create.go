package club

import (
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

//创建机器人
//参数:
//	1.guildInfo 工会信息
//	2.robotUsers 所有机器人
//	3.normalUsers 正常用户
//返回:
//	1.新创建的用户
//  2.错误信息
func CreateRobot(guildInfo model.Guild, robotUsers []model.User, normalUsers []model.User) (model.User, error) {
	//机器人姓名不能重复
	var existsNames []string
	for _, u := range robotUsers {
		existsNames = append(existsNames, u.UserName)
	}

	userName := getRobotNameByRand(existsNames...)
	userHeadIcon := getRobotIconByRand()
	userLikeCount := getLikeNumByRand()
	userLevel := getLevelByRand(guildInfo, normalUsers)

	//create
	u, err := service.CreateGuildRobot(userName, userHeadIcon, userLikeCount, userLevel)
	if err != nil {
		return u, err
	}
	return u, err
}
