package service

import (
	"math/rand"
	"testing"
	"time"
)

func Test_SendChatMessage(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	/*
		u, err := CreateGuildRobot("Levin666", "4", 19, 97)
		if err != nil {
			t.Fatal(err)
		}
		if err := JoinToGuild(GuildID, u.UserID); err != nil {
			t.Fatal(err)
		}
	*/
	u, err := GetUserInfoByUserID(130056000009)
	if err != nil {
		t.Fatal(err)
	}
	if err := SendChatMessage(GuildID, u, "你好"); err != nil {
		t.Fatal(err)
	}

	if err := SendJoinRoomMessage(GuildID, u); err != nil {
		t.Fatal(err)
	}
}
