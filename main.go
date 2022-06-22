package main

import (
	"flag"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/casual-server-lib/redis"
	"github.com/joycastle/matching-story-robot-service/confmanager"
)

func main() {
	//run params
	runEnv := flag.String("env", "dev", "dev(本机开发环境), pre(预发布环境), prod(线上环境)")
	flag.Parse()

	configFile := "./conf/" + *runEnv + ".yaml"
	if err := config.InitConfig(configFile); err != nil {
		panic(err)
	}

	//print configs
	log.Infof("server-lib config filePath: %s", configFile)
	log.Infof("server-lib log config: %v", config.Logs)
	log.Infof("server-lib redis config: %v", config.Redis)
	log.Infof("server-lib mysql config: %v", config.Mysql)

	//init logs
	log.InitLogs(config.Logs)

	//init mysql
	if err := mysql.InitMysql(config.Mysql); err != nil {
		panic(err)
	}

	// init redis
	redis.InitRedis(config.Redis)

	//confmanager
	if err := confmanager.GetConfManagerVer().LoadCsv("confmanager/template"); err != nil {
		panic(err)
	}
}
