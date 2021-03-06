package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/joycastle/matching-story-robot-service/confmanager"
)

func TestMain(m *testing.M) {
	//init configmanager
	if err := confmanager.GetConfManagerVer().LoadCsv("../../confmanager/template"); err != nil {
		panic(err)
	}
	if err := Startup(); err != nil {
		panic(err)
	}
	m.Run()
}

func TestRobotTeamConfig(t *testing.T) {
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

	fmt.Println(GetHelpTalkTimeGapByRand())
}

func TestRobotTeamChat(t *testing.T) {
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

func TestGetActiveTimeByRand(t *testing.T) {
	fmt.Println("GetActiveTimeByRand", GetActiveTimeByRand())
}

func TestGetRobotKictTimeRange(t *testing.T) {
	GetRobotKictTimeRange()
	fmt.Println(GetRobotKictNum())
}

func TestGetFirstActionTimeByRand(t *testing.T) {
	f, err := GetFirstActionTimeByRand(10090)
	fmt.Println("GetFirstActionTimeByRand", f, err)
}

func TestGetRule2TargetByRand(t *testing.T) {
	fmt.Println("GetRule2TargetByRand")
	for i := 1001; i < 1010; i++ {
		fmt.Println(GetRule2TargetByRand(i))
	}
}

func TestAAA(t *testing.T) {
	fmt.Println(GetRobotActionIDByRand(300, 2))
}
