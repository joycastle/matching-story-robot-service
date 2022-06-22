package club

import (
	"testing"

	"github.com/joycastle/matching-story-robot-service/model"
	"github.com/joycastle/matching-story-robot-service/service"
)

var GuildID int64 = 72719317539487744

func Test_CreateRobot(t *testing.T) {
	guildInfo, err := service.GetGuildInfoByGuildID(GuildID)
	if err != nil {
		t.Fatal(err)
	}

	userDistributions, err := service.GetGuildUserTypeDistribution(GuildID)
	if err != nil {
		t.Fatal(err)
	}

	var robotUsers []model.User
	var normalUsers []model.User
	if v, ok := userDistributions[service.USERTYPE_CLUB_ROBOT_SERVICE]; ok {
		robotUsers = v
	}
	if v, ok := userDistributions[service.USERTYPE_NORMAL]; ok {
		normalUsers = v
	}

	robotUser, err := CreateRobot(guildInfo, robotUsers, normalUsers)
	if err != nil {
		t.Fatal(err)
	} else if robotUser.UserType != service.USERTYPE_CLUB_ROBOT_SERVICE {
		t.Fatal(robotUser)
	}
}
