package create

import (
	"fmt"

	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/model"
)

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
