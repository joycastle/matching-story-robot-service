package service

import (
	"time"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/casual-server-lib/redis"
	"github.com/joycastle/matching-story-robot-service/confmanager"
)

func init() {
	//init configmanager
	confMgr := confmanager.GetConfManagerVer()
	err := confMgr.LoadCsv("../confmanager/template")
	if err != nil {
		panic(err)
	}

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
			SlowSqlTime: 0,
			SlowLogger:  "slow",
			ErrLogger:   "error",
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
			SlowSqlTime: 0,
			SlowLogger:  "slow",
			ErrLogger:   "error",
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
}
