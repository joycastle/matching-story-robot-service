package create

import (
	"sync"

	"github.com/joycastle/casual-server-lib/faketime"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/library"
	"github.com/joycastle/matching-story-robot-service/lib"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	kickRecordMap map[int64]int64 = make(map[int64]int64, 5000)
	kickRecordMu  *sync.RWMutex   = new(sync.RWMutex)
)

func IsKickingState(gid int64) bool {
	kickRecordMu.RLock()
	defer kickRecordMu.RUnlock()

	if expire, ok := kickRecordMap[gid]; ok {
		if faketime.Now().Unix()-expire < 0 {
			return true
		} else {
			delete(kickRecordMap, gid)
		}
	}

	return false
}

func clearKickingState(gid int64) int {
	kickRecordMu.Lock()
	defer kickRecordMu.Unlock()

	if _, ok := kickRecordMap[gid]; ok {
		delete(kickRecordMap, gid)
	}
	length := len(kickRecordMap)
	return length
}

func markKickingState(gid int64) int {
	kickRecordMu.Lock()
	defer kickRecordMu.Unlock()

	kickRecordMap[gid], _ = kickRobotTimeHandler(nil)
	length := len(kickRecordMap)

	return length
}

func kickRobotTimeHandler(*library.Job) (int64, error) {
	return faketime.Now().Unix() + config.GetRobotKictTimeRange(), nil
}

func kickRobotLogicHandler(job *library.Job) *lib.LogStructuredJson {
	info := lib.NewLogStructed()
	kickRecordLen := 0
	kickRecordLen = clearKickingState(job.GuildID)

	//获取机器人分布
	uids, err := service.GetGuildUserIds(job.GuildID)
	if err != nil {
		return info.Failed().Step(10).Err(err)
	}
	//是否满员
	if len(uids) >= 30 {
		return info.Failed().Step(1).ErrString("guild members is full")
	}
	if len(uids) == 0 {
		return info.Failed().Step(2).ErrString("guild members is empty")
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
			if (faketime.Now().Unix()-u.LoginTime/1000)/86400 < int64(avtiveDays) {
				normalUserActiveNum++
			}
		}
	}

	lowerLimit := config.GetRobotKickNumLimit()
	if robotNum <= lowerLimit {
		return info.Failed().Step(3).ErrString("guild robot num is lower csv limit").Set("robotNum", robotNum, "csv:robot_num_minlimit", lowerLimit)
	}

	avtiveUserNum := config.GetRobotKickActiveUsersNum()
	if normalUserActiveNum <= avtiveUserNum {
		return info.Failed().Step(4).ErrString("guild robot num is lower csv limit").Set("avtiveUserNum", normalUserActiveNum, "csv:robot_kick_active_username", avtiveUserNum, "csv:robot_kick_level", avtiveDays)
	}

	kickRecordLen = markKickingState(job.GuildID)

	//发送退出工会记录
	respone, err := service.SendLeaveGuildRPC("", robotMaxLevelUid, job.GuildID)
	if err != nil {
		return info.Failed().Step(5).Err(err).Set("respone", respone)
	}

	return info.Success().Set("leaveUserId", robotMaxLevelUid, "kickRecordLen", kickRecordLen)
}
