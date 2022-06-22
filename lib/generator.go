package lib

import "github.com/joycastle/casual-server-lib/redis"

//生成userID
func GenerateUserID() (int64, error) {
	redisKey := "MaxUserIDIncrNumKey"
	index, err := redis.Incr("default", redisKey)
	if err != nil {
		return 0, err
	}

	zoneId := 0
	gUserPreZoneBase := 1000
	gUserPreLanguageBase := 10
	language := 0
	deviceType := 9

	return int64(1000000*index) + int64(zoneId*gUserPreZoneBase) + int64(language*gUserPreLanguageBase) + int64(deviceType), nil
}
