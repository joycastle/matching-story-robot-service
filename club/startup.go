package club

import (
	"github.com/joycastle/matching-story-robot-service/club/action"
	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/club/create"
)

func StartupServiceRobot() {
	config.Startup()
	create.Startup()
	action.Startup()
}
