package service

import (
	"testing"
)

var GuildID int64 = 72719317539487744

func Test_GetClubUserTypeDistribution(t *testing.T) {
	disMap, err := GetClubUserTypeDistribution(GuildID)
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

func Test_JoinToClub(t *testing.T) {
	disMap, err := GetClubUserTypeDistribution(GuildID)
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

	u, err := CreateClubRobot(name, icon, likeCnt, level)
	if err != nil {
		t.Fatal(err)
	}

	if err := JoinToClub(GuildID, u.UserID); err != nil {
		t.Fatal(err)
	}

	disMap2, err := GetClubUserTypeDistribution(GuildID)
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

	if len(disMap[USERTYPE_CLUB_ROBOT]) != len(disMap2[USERTYPE_CLUB_ROBOT])-1 {
		t.Fatal("3")
	}
}
