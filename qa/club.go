package qa

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type QaClubAction struct {
	Time     string
	OldLevel int
	NewLevel int
	OldTimes int
	NewTimes int
	Msg      string
}

var clubActionMap map[int64]map[int64][]QaClubAction = make(map[int64]map[int64][]QaClubAction)
var clubActionMapMu *sync.Mutex = new(sync.Mutex)

func init() {
	qaDebug.AddInt("guild_id", "工会id")
}

func GetGuildIDString() string {
	return qaDebug.Get("guild_id")
}

func GetGuildIDMap() map[int64]struct{} {
	list := strings.Split(qaDebug.Get("guild_id"), ",")
	m := make(map[int64]struct{})
	for _, v := range list {
		i64, _ := strconv.ParseInt(v, 10, 64)
		m[i64] = struct{}{}
	}
	return m
}

func GetGuildIDSlice() []int64 {
	var ret []int64
	list := strings.Split(qaDebug.Get("guild_id"), ",")
	for _, v := range list {
		i64, _ := strconv.ParseInt(v, 10, 64)
		ret = append(ret, i64)
	}
	return ret
}

func IsExistsInQaDebug(guildID int64) bool {
	m := GetGuildIDMap()
	_, ok := m[guildID]
	return ok
}

func AddGuildAction(guildiD, uid int64, oldLevel, newLevel, oldTimes, newTimes int) {
	clubActionMapMu.Lock()
	defer clubActionMapMu.Unlock()

	if _, ok := clubActionMap[guildiD]; !ok {
		clubActionMap[guildiD] = make(map[int64][]QaClubAction)
	}

	if _, ok := clubActionMap[guildiD][uid]; !ok {
		clubActionMap[guildiD][uid] = []QaClubAction{}
	}

	var tmp QaClubAction
	tmp.Time = time.Now().Format("2006-01-02 15:04:05")
	tmp.OldLevel = oldLevel
	tmp.NewLevel = newLevel
	tmp.OldTimes = oldTimes
	tmp.NewTimes = newTimes

	clubActionMap[guildiD][uid] = append(clubActionMap[guildiD][uid], tmp)
}

func AddGuildActionError(guildiD, uid int64, msg string) {
	clubActionMapMu.Lock()
	defer clubActionMapMu.Unlock()

	if _, ok := clubActionMap[guildiD]; !ok {
		clubActionMap[guildiD] = make(map[int64][]QaClubAction)
	}

	if _, ok := clubActionMap[guildiD][uid]; !ok {
		clubActionMap[guildiD][uid] = []QaClubAction{}
	}

	var tmp QaClubAction
	tmp.Time = time.Now().Format("2006-01-02 15:04:05")
	tmp.Msg = msg

	clubActionMap[guildiD][uid] = append(clubActionMap[guildiD][uid], tmp)
}

func GetGuildActionReport(guildiD int64) string {
	m := make(map[int64][]QaClubAction)
	clubActionMapMu.Lock()
	if _, ok := clubActionMap[guildiD]; !ok {
		return fmt.Sprintf("当前工会未记录:%d", guildiD)
	}

	for k, vs := range clubActionMap[guildiD] {
		if _, ok := m[k]; !ok {
			m[k] = []QaClubAction{}
		}
		for _, v := range vs {
			m[k] = append(m[k], v)
		}
	}
	clubActionMapMu.Unlock()

	ret := ""
	for k, vs := range m {
		for _, v := range vs {
			ret += fmt.Sprintf("guild_id:%d, uid:%d, msg:%s\n", guildiD, k, v)
		}
	}
	return ret
}
