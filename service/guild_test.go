package service

import "testing"

func Test_GetGuildInfoByGuildID(t *testing.T) {
	if g, err := GetGuildInfoByGuildID(GuildID); err != nil || g.ID != GuildID {
		t.Fatal(err, g)
	}
}
