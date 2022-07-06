package service

import (
	"testing"
	"time"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/casual-server-lib/redis"
)

var isConfiged bool = false

func TestGetGuildRequestInfosWithFiledsByGuildIDs(t *testing.T) {
	testinit()
	_, err := GetGuildRequestInfosWithFiledsByGuildIDs([]int64{9068658676465664, 9187840659292160}, []string{"done"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGuildResponeInfosWithFiledsByHelpIDs(t *testing.T) {
	testinit()
	_, err := GetGuildResponeInfosWithFiledsByHelpIDs([]int64{9207422417633280}, []string{"responder_id", "time"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUserInfosWithField(t *testing.T) {
	testinit()
	_, err := GetUserInfosWithField([]int64{48000060, 213000100, 3675}, []string{"user_name"})
	if err != nil {
		t.Fatal(err)
	}
}

func testinit() {
	if !isConfiged {
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
		config.Grpc["chat"] = "http://127.0.0.1:8081"

		redis.InitRedis(redisConfigs)

		isConfiged = true
	}
}
