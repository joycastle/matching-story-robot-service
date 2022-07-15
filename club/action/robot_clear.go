package action

import (
	"sync"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/cache"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

func clearRobotActiveTimeHandler(job *GuildJob) (int64, *Result) {
	rnd := config.GetRobotKictTimeRange()
	return faketime.Now().Unix() + rnd, ActionSuccess().Detail("rnd", rnd)
}

//记录机器人清退状态
var (
	kickStatusMap map[int64]int64 = make(map[int64]int64)
	kickStatusMu  *sync.RWMutex   = new(sync.RWMutex)
)

func MarkGuildKickStatus(guildId int64) {
	kickStatusMu.Lock()
	defer kickStatusMu.Unlock()
	kickStatusMap[guildId] = faketime.Now().Unix() + config.GetRobotKictTimeRange()
}

func clearRobotDispatchHandler(job *GuildJob) *Result {
	deleted, err := cache.IsGuildDelete(job.GuildID)
	if err != nil {
		return ErrorText(100).Detail("guild", job.GuildID, err)
	}
	if deleted {
		return ErrorText(5000).Detail("id", job.GuildID)
	}

	//获取机器人分布
	uids, err := service.GetGuildUserIds(job.GuildID)
	if err != nil {
		return ErrorText(100).Detail("table:guild", "guild_id", job.GuildID, err)
	}
	//是否满员
	if len(uids) >= 30 {
		return ErrorText(5001).Detail("id", job.GuildID)
	}
	if len(uids) == 0 {
		return ErrorText(5002).Detail("id", job.GuildID)
	}
	//判断机器人数量下线
	robotNum := 0
	normalUserActiveNum := 0
	avtiveDays := config.GetRobotKickActiveDays()
	users, err := service.GetUserInfosWithField(uids, []string{"login_time"})
	robotMaxLevelUid := int64(0)
	robotMaxLevel := 0
	for _, u := range users {
		if u.UserType == service.USERTYPE_CLUB_ROBOT_SERVICE {
			robotNum++
			if u.UserLevel > robotMaxLevel {
				robotMaxLevel = u.UserLevel
				robotMaxLevelUid = u.UserID
			}
		} else if u.UserType == service.USERTYPE_NORMAL {
			if (faketime.Now().Unix()-u.LoginTime/1000)/86400 > int64(avtiveDays) {
				normalUserActiveNum++
			}
		}
	}

	lowerLimit := config.GetRobotKickNumLimit()
	if robotNum <= lowerLimit {
		return ErrorText(5003).Detail("id", job.GuildID, "robot num", robotNum, "csv:robot_num_minlimit", lowerLimit)
	}

	avtiveUserNum := config.GetRobotKickActiveUsersNum()
	if normalUserActiveNum <= avtiveUserNum {
		return ErrorText(5004).Detail("id", job.GuildID, "avtive user num", normalUserActiveNum, "csv:robot_kick_active_username", avtiveUserNum, "csv:robot_kick_level", avtiveDays)
	}

	//发送退出工会记录
	CreateLeaveGuildJob(robotMaxLevelUid, job.GuildID)

	//标记清退状态
	MarkGuildKickStatus(job.GuildID)

	return ActionSuccess().Detail("leave user_id", robotMaxLevelUid)
}
