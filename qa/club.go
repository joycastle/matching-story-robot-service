package qa

import (
	"strconv"
	"strings"
)

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
