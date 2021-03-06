package service

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/casual-server-lib/redis"
)

func TestMain(m *testing.M) {
	//init mysql
	mysqlConfigs := map[string]mysql.MysqlConf{
		"default-master": mysql.MysqlConf{
			Addr:        "127.0.0.1",
			Username:    "root",
			Password:    "123456",
			Database:    "db_game",
			Options:     "charset=utf8mb4&parseTime=True",
			MaxIdle:     16,
			MaxOpen:     128,
			MaxLifeTime: time.Second * 300,
			SlowSqlTime: 1,
			SlowLogger:  "slow",
			StatLogger:  "stat",
		},

		"default-slave": mysql.MysqlConf{
			Addr:        "127.0.0.1",
			Username:    "root",
			Password:    "123456",
			Database:    "db_game",
			Options:     "charset=utf8mb4&parseTime=True",
			MaxIdle:     16,
			MaxOpen:     128,
			MaxLifeTime: time.Second * 300,
			SlowSqlTime: 1,
			SlowLogger:  "slow",
			StatLogger:  "stat",
		},
	}

	if err := mysql.InitMysql(mysqlConfigs); err != nil {
		panic(err)
	}

	//init redis
	redisConfigs := map[string]redis.RedisConf{
		"default": redis.RedisConf{
			Addr:           "127.0.0.1:6379,127.0.0.1:6379,127.0.0.1:6379",
			Password:       "123456",
			MaxActive:      32,
			MaxIdle:        16,
			IdleTimeout:    time.Second * 1800,
			ConnectTimeout: time.Second * 10,
			ReadTimeout:    time.Second * 2,
			WriteTimeout:   time.Second * 2,
			TestInterval:   time.Second * 300,
		},
	}

	//init grpc
	config.Grpc = make(map[string]string)
	config.Grpc["default"] = "http://127.0.0.1:3002"

	redis.InitRedis(redisConfigs)

	m.Run()
}

func UpdateRobotRule2Config(userId int64, m map[int]int) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := UpdateRobotRule2ConfigByUid(userId, string(b)); err != nil {
		return err
	}
	return nil
}

func TestRpc(t *testing.T) {
	fmt.Println(MergeFiledsKV("name", "levin", "age", 11, "122", 12))
}

/*
func TestRpc(t *testing.T) {
	m := make(map[int]int)
	m[0] = 0
	m[1] = 1
	m[2] = 3
	UpdateRobotRule2Config(136831000009, m)

	mm := make(map[int]int)
	u, _ := GetRobotForGuild(136831000009)
	fmt.Println(u.Name, []byte(u.Name))
	json.Unmarshal([]byte(u.Name), &mm)
	fmt.Println(mm)
}

func TestRpc(t *testing.T) {
	return
	if _, err := SendUpdateScoreRPC("43f6a40db954b4913c47ed60fa665bf9", 136833000009, 2); err != nil {
		t.Fatal(err)
	}
}

func TestRpcLeaveGuild(t *testing.T) {
	fmt.Println(SendLeaveGuildRPC("", 136845000009, 131575164054798336))
}

func TestGetGuildRequestInfosWithFiledsByGuildIDs(t *testing.T) {
	_, err := GetGuildRequestInfosWithFiledsByGuildIDs([]int64{9068658676465664, 9187840659292160}, []string{"done"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGuildResponeInfosWithFiledsByHelpIDs(t *testing.T) {
	_, err := GetGuildResponeInfosWithFiledsByHelpIDs([]int64{9207422417633280}, []string{"responder_id", "time"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserInfosWithField(t *testing.T) {
	_, err := GetUserInfosWithField([]int64{48000060, 213000100, 3675}, []string{"user_name"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateGuildRobotUserRPC(t *testing.T) {
	fmt.Println(CreateGuildRobotUserRPC("111111", "2", 1999, 88))
}*/
