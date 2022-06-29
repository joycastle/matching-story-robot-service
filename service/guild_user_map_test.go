package service

import (
	"testing"
)

var GuildID int64 = 72719317539487744

func Test_GetGuildUserTypeDistribution(t *testing.T) {
	disMap, err := GetGuildUserTypeDistribution(GuildID)
	if err != nil {
		t.Fatal(err)
	}

	if len(disMap) < 2 {
		t.Fatal("Maybe Error, global length")
	}

	for k, v := range disMap {
		if len(v) <= 0 {
			t.Fatal("Maybe Error internal length", k)
		}
	}
}

func Test_JoinToGuild(t *testing.T) {
	disMap, err := GetGuildUserTypeDistribution(GuildID)
	if err != nil {
		t.Fatal(err)
	}

	totalNum := 0
	for _, v := range disMap {
		totalNum = totalNum + len(v)
	}

	name := "TestHello"
	icon := "8"
	likeCnt := 88
	level := 99

	u, err := CreateGuildRobotUser(name, icon, likeCnt, level)
	if err != nil {
		t.Fatal(err)
	}

	if err := JoinToGuild(GuildID, u.UserID); err != nil {
		t.Fatal(err)
	}

	disMap2, err := GetGuildUserTypeDistribution(GuildID)
	if err != nil {
		t.Fatal(err)
	}

	totalNum2 := 0
	for _, v := range disMap2 {
		totalNum2 = totalNum2 + len(v)
	}

	if totalNum != totalNum2-1 {
		t.Fatal("1")
	}

	if len(disMap[USERTYPE_NORMAL]) != len(disMap2[USERTYPE_NORMAL]) {
		t.Fatal("2")
	}

	if len(disMap[USERTYPE_CLUB_ROBOT_SERVICE]) != len(disMap2[USERTYPE_CLUB_ROBOT_SERVICE])-1 {
		t.Fatal("3")
	}
}

func TestBatchGetGuildUserIdsByGuildIDs(t *testing.T) {
	m, err := BatchGetGuildUserIdsByGuildIDs([]int64{9068658676465664, 9194785835319296})
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Fatal("")
	}
}

func TestGetAllGuildUserMapInfos(t *testing.T) {
}
