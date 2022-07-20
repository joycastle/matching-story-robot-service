package create

import (
	"fmt"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/action"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

func createRobotTimeHandler(*library.Job) (int64, error) {
	return faketime.Now().Unix() + config.GetActiveTimeByRand(), nil
}

func createRobotLogicHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()

	if IsKickingState(job.GuildID) {
		return info.Failed().Step(10).ErrString("guild is kicking state")
	}

	//1.从DB获取工会信息
	guildInfo, err := service.GetGuildInfoByGuildID(job.GuildID)
	if err != nil {
		return info.Failed().Step(11).Err(err)
	}

	//2.获取机器人和真实用户分布
	userDistributions, err := service.GetGuildUserTypeDistribution(job.GuildID)
	if err != nil {
		return info.Failed().Step(1).Err(err)
	}

	var robotUsers []model.User
	var normalUsers []model.User
	if v, ok := userDistributions[service.USERTYPE_CLUB_ROBOT_SERVICE]; ok {
		robotUsers = v
	}
	if v, ok := userDistributions[service.USERTYPE_NORMAL]; ok {
		normalUsers = v
	}

	//case 全是机器人则删除工会
	if len(normalUsers) == 0 {
		return info.Failed().Step(6).ErrString("will delete guild").Set("normalUsers", len(normalUsers), "robotUsers", len(robotUsers))
	}

	//3.判断是否满足创建机器人的条件
	reason, ok := IsTheGuildCanUsingRobot(guildInfo, robotUsers, normalUsers)
	if !ok {
		return info.Failed().Step(3).Set("reason", reason)
	}

	//4.创建机器人
	//4.1获取随机创建机器人数量
	newNum := config.GetGenerateRobotNumByRand()
	robotMaxLimitNum := config.GetRobotMaxLimitNum()
	//4.2边界检查
	if len(robotUsers)+newNum > robotMaxLimitNum {
		newNum = robotMaxLimitNum - len(robotUsers)
	}

	//newNum一定>0

	//4.3创建机器人
	var newRobots []model.User
	var newRobotsUid map[int64]int = make(map[int64]int)
	isCreateOk := true
	for i := 0; i < newNum; i++ {
		if rbtUser, err := CreateRobot(guildInfo, robotUsers, normalUsers); err != nil {
			info.Failed().Step(4).Err(err)
			isCreateOk = false
			break
		} else {
			//名称过滤用
			robotUsers = append(robotUsers, rbtUser)
			//日志打印和加入工会用
			newRobots = append(newRobots, rbtUser)
			newRobotsUid[rbtUser.UserID] = rbtUser.UserLevel
		}
	}
	if !isCreateOk {
		//只要有一个失败，则不执行加入工会操作
		return info
	}

	//5.加入工会
	for _, ru := range newRobots {
		if _, err := service.SendJoinToGuildRPC(ru.AccountID, ru.UserID, job.GuildID); err != nil {
			return info.Failed().Step(5).Err(err).Set("join_user_id", ru.UserID)
		}
		action.CreateFirstInJob(ru.UserID, job.GuildID)
	}

	return info.Success().Set("newNum", newNum, "createNum", len(newRobots), "levelInfo", newRobotsUid)
}

//创建机器人
//参数:
//  1.guildInfo 工会信息
//  2.robotUsers 所有机器人
//  3.normalUsers 正常用户
//返回:
//  1.新创建的用户
//  2.错误信息
func CreateRobot(guildInfo model.Guild, robotUsers []model.User, normalUsers []model.User) (model.User, error) {
	//机器人姓名不能重复
	var existsNames []string
	for _, u := range robotUsers {
		existsNames = append(existsNames, u.UserName)
	}

	userName := config.GetRobotNameByRand(existsNames...)
	userHeadIcon := config.GetRobotIconByRand()
	userLikeCount := config.GetLikeNumByRand()
	userLevel := config.GetLevelByRand(guildInfo, normalUsers)

	//create user
	u, err := service.CreateGuildRobotUserRPC(userName, userHeadIcon, userLikeCount, userLevel)
	if err != nil {
		return u, err
	}

	//create robot config
	robotUserType := config.GetRobotTypeByRand(userLevel)
	robotAction := config.GetRobotActionIDByRand(userLevel, robotUserType)
	robotRule1SleepTs, err := config.GetRule1TargetByRand(int(robotAction))
	if err != nil {
		return u, err
	}
	_, err = service.CreateRobotForGuild(u.UserID, robotUserType, robotAction, robotRule1SleepTs)
	if err != nil {
		return u, err
	}

	return u, nil
}

//判断当前工会是否满足开启机器人的条件
//参数:
//  1.guildInfo 工会信息
//  2.robotUsers 所有机器人
//  3.normalUsers 正常用户
//返回:
//  reason 不满足的原因
//  true 满足，false 不满足
func IsTheGuildCanUsingRobot(guildInfo model.Guild, robotUsers []model.User, normalUsers []model.User) (string, bool) {
	if guildInfo.DeletedAt.Valid {
		return "The club has been deleted", false
	}

	if guildInfo.IsOpen != 2 {
		return "The club not allowed to join", false
	}

	robotMaxLimitNum := config.GetRobotMaxLimitNum()
	if len(robotUsers) >= robotMaxLimitNum {
		return fmt.Sprintf("The robot number has reached the maxmum limit of %d", robotMaxLimitNum), false
	}

	normalUserLimitNum := config.GetNormalUserNum()
	if len(normalUsers) >= normalUserLimitNum {
		return fmt.Sprintf("The normal number has exceeded the open limit of %d", normalUserLimitNum), false
	}

	return "", true
}
