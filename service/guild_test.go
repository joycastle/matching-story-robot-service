package service

import (
	"fmt"
	"testing"
)

func Test_GetGuildInfoByGuildID(t *testing.T) {
	if g, err := GetGuildInfoByGuildID(GuildID); err != nil || g.ID != GuildID {
		t.Fatal(err, g)
	}
}

func TestGetAllGuildInfoFromDBWithSliceArray(t *testing.T) {
	list, err := GetAllGuildInfoFromDB()
	if err != nil {
		t.Fatal(err)
	}

	listArray, err := GetAllGuildInfoFromDBWithSliceArray()
	if err != nil {
		t.Fatal(err)
	}

	total := 0
	for _, v := range listArray {
		total = total + len(v)
	}

	if len(list) != total {
		t.Fatal("len(list) != total")
	}
}

func TestGetGuildDeleteInfos(t *testing.T) {
	if ret, err := GetGuildDeleteInfos([]int64{9274660155817984}); err != nil || len(ret) != 1 {
		t.Fatal(err)
	} else {
		fmt.Println(ret)
	}

	if ret, err := GetGuildDeleteInfos([]int64{}); err != nil || len(ret) != 0 {
		t.Fatal(err, ret)
	}

	if ret, err := GetGuildDeleteInfos([]int64{96000060}); err != nil || len(ret) != 0 {
		t.Fatal(err, ret)
	} else {
		fmt.Println(ret, err)
	}
}
