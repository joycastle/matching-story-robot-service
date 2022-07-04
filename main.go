package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/joycastle/casual-server-lib/config"
	"github.com/joycastle/casual-server-lib/log"
	"github.com/joycastle/casual-server-lib/mysql"
	"github.com/joycastle/casual-server-lib/redis"
	"github.com/joycastle/matching-story-robot-service/club"
	"github.com/joycastle/matching-story-robot-service/confmanager"
)

func main() {
	//run params
	runEnv := flag.String("env", "dev", "dev(本机开发环境), pre(预发布环境), prod(线上环境), other_env(其他开发配置名称)")
	flag.Parse()

	if *runEnv == "" {
		panic("环境选择错误")
	}

	configFileEnv := filepath.Join("./conf/", *runEnv)
	configFileCommon := filepath.Join("./conf/common")
	if err := config.InitConfig(configFileEnv, configFileCommon); err != nil {
		panic(err)
	}

	//print configs
	log.Infof("server-lib config env filePath: %s", configFileEnv)
	log.Infof("server-lib config common filePath: %s", configFileCommon)
	log.Infof("server-lib log config: %v", config.Logs)
	log.Infof("server-lib redis config: %v", config.Redis)
	log.Infof("server-lib mysql config: %v", config.Mysql)
	log.Infof("server-lib grpc config: %v", config.Grpc)

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

	//启动服务型机器人
	club.StartupServiceRobot()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
}
