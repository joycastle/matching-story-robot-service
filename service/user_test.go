package service

import (
	"testing"
)

func Test_CreateGuildRobot_GetOne_BatchGet(t *testing.T) {
	name := "TestHello"
	icon := "8"
	likeCnt := 88
	level := 99

	u, err := CreateGuildRobot(name, icon, likeCnt, level)
	if err != nil {
		t.Fatal(err)
	}

	//get
	if _, err := GetUserInfoByUserID(u.UserID + 1); err != ErrRecordNotFound {
		t.Fatal(err)
	}

	if u3, err := GetUserInfoByUserID(u.UserID); err != nil {
		t.Fatal(err)
	} else {
		if u3.UserName != name || u3.UserHeadIcon != icon || u3.UserLikeCount != uint(likeCnt) || u3.UserLevel != level {
			t.Fatal("CreateGuildRobot Error 1")
		}

		if u3.DeviceType != DEVICE_TYPE_CLUB_ROBOT || u3.UserType != USERTYPE_CLUB_ROBOT_SERVICE || u3.UserCountryData != COUNTRY_CN {
			t.Fatal("CreateGuildRobot Error 2")
		}
	}

	//batch get
	if u4s, err := BatchGetUserInfoByUserID([]int64{u.UserID + 1}); len(u4s) > 0 {
		t.Fatal(err, u4s)
	}

	var uids []int64
	var uidMap map[int64]struct{} = make(map[int64]struct{})
	if u, err := CreateGuildRobot(name, icon, likeCnt, level); err != nil {
		t.Fatal(err)
	} else {
		uids = append(uids, u.UserID)
		uidMap[u.UserID] = struct{}{}
	}
	if u, err := CreateGuildRobot(name, icon, likeCnt, level); err != nil {
		t.Fatal(err)
	} else {
		uids = append(uids, u.UserID)
		uidMap[u.UserID] = struct{}{}
	}
	if u, err := CreateGuildRobot(name, icon, likeCnt, level); err != nil {
		t.Fatal(err)
	} else {
		uids = append(uids, u.UserID)
		uidMap[u.UserID] = struct{}{}
	}

	if u5s, err := BatchGetUserInfoByUserID(uids); len(u5s) != len(uids) || err != nil {
		t.Fatal(err, u5s)
	} else {
		for _, v := range u5s {
			if _, ok := uidMap[v.UserID]; !ok {
				t.Fatal("BatchGetUserInfoByUserID result not match uids")
			}
		}
	}
}
