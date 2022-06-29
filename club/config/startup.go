package config

import (
	"math/rand"
	"time"
)

func init() {
	go func() {
		for {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Second)
		}
	}()
}

func Startup() error {
	if err := ReadRobotTeamConfigFromConfManager(); err != nil {
		return err
	}

	if err := ReadRobotTeamChatFromConfManager(); err != nil {
		return err
	}

	if err := ReadRobotTeamInitialFromConfManager(); err != nil {
		return err
	}

	if err := ReadRobotTeamFromConfManager(); err != nil {
		return err
	}

	return nil
}
