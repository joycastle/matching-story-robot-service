package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/joycastle/matching-story-robot-service/confmanager"
)

func ConfigInit() {
	//init configmanager
	if err := confmanager.GetConfManagerVer().LoadCsv("../../confmanager/template"); err != nil {
		panic(err)
	}

	if err := Startup(); err != nil {
		panic(err)
	}
}

func TestRobotTeamConfig(t *testing.T) {
	ConfigInit()

	if GetRobotMaxLimitNum() <= 0 || GetNormalUserNum() <= 0 {
		t.Fatal("")
	}

	ok := false
	for i := 0; i < 200; i++ {
		num := GetGenerateRobotNumByRand()
		if num > 2 || num < 1 {
			t.Fatal("getGenerateRobotNumByRand", num)
		}
		if num == 2 {
			ok = true
		}
		time.Sleep(1000)
	}
	if !ok {
		t.Fatal("getGenerateRobotNumByRand error")
	}

	ok2 := false
	for i := 0; i < 200; i++ {
		num := GetLikeNumByRand()
		if num > 15 || num < 5 {
			t.Fatal("getLikeNumByRand")
		}
		if num == 15 {
			ok2 = true
		}
		time.Sleep(1000)
	}
	if !ok2 {
		t.Fatal("getLikeNumByRand error")
	}

	ok3 := false
	for i := 0; i < 200; i++ {
		num := GetJoinTalkTimeGapByRand()
		if num > 120 || num < 60 {
			t.Fatal("GetJoinTalkTimeGapByRand")
		}
		if num == 120 || num == 60 {
			ok3 = true
		}
		time.Sleep(1000)
	}
	if !ok3 {
		t.Fatal("GetJoinTalkTimeGapByRand error")
	}
}

func TestRobotTeamChat(t *testing.T) {
	ConfigInit()
	if len(GetChatMsgByRand(1)) == 0 {
		t.Fatal("GetChatMsgByRand index 1 error")
	}
	if len(GetChatMsgByRand(2)) == 0 {
		t.Fatal("GetChatMsgByRand index 2 error")
	}
}

func TestRobotTeamInitial(t *testing.T) {
	m := make(map[int]int, 3)
	for i := 0; i < 100; i++ {
		index := getWeightIndex([]int{5, 1, 1})
		time.Sleep(2000)
		if _, ok := m[index]; !ok {
			m[index] = 0
		}
		m[index]++
	}

	fmt.Println(m)

	fmt.Println(GetRobotActionIDByRand(133, 1))
	fmt.Println(GetRobotTypeByRand(133))

}

func TestRobotTeam(t *testing.T) {
	fmt.Println(GetRobotActiveDaysByActionID(1009))
}
