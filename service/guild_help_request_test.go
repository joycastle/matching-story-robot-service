package service

import (
	"testing"
	"time"

	"gorm.io/gorm"
)

func Test_GetMinGuildHelpRequestIDByTimeAfter(t *testing.T) {
	curTime := time.Now().Unix() + 86400*1000
	id, err := GetMinGuildHelpRequestIDByTimeAfter(curTime)
	if err != gorm.ErrRecordNotFound {
		t.Fatal(err, id)
	}

	curTime = time.Now().Unix() - 86400*10
	id, err = GetMinGuildHelpRequestIDByTimeAfter(curTime)
	if (err != nil || id <= 0) && err != gorm.ErrRecordNotFound {
		t.Fatal(err, id)
	}

	info, err := GetGuildHelpRequestInfoByID(id)
	if err != nil {
		t.Fatal(err)
	}

	if info.Time < curTime {
		t.Fatal("logic error")
	}

	maxID, err := GetMaxGuildHelpRequestID()
	if err != nil {
		t.Fatal(err, maxID)
	}

	count, err := BatchGetGuildHelpRequestCountByAfterID(-1)
	if err != nil {
		t.Fatal(err)
	}

	if maxID <= 0 && count > 0 {
		t.Fatal("logic error", maxID, count)
	}
}

func Test_BatchGetGuildHelpRequestInfoByAfterID(t *testing.T) {
	curTime := time.Now().Unix() - 86400*100
	id, err := GetMinGuildHelpRequestIDByTimeAfter(curTime)
	if err == gorm.ErrRecordNotFound {
		t.Fatal(err, id)
	}
	list, err := BatchGetGuildHelpRequestInfoByAfterID(id)
	if err != nil {
		t.Fatal(err)
	}
	count, err := BatchGetGuildHelpRequestCountByAfterID(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != count {
		t.Fatal("logic error")
	}
}

func Test_GetGuildHelpRequestNotComplete(t *testing.T) {
	/*
		if list, err := GetGuildHelpRequestNotComplete(); err != nil {
			if err != gorm.ErrRecordNotFound {
				t.Fatal(err)
			}
		} else {
			for _, one := range list {
				if one.Total == one.Count {
					t.Fatal(list[0])
				}
			}
		}*/
}

func Test_UpdateGuildHelpRequestCountByID(t *testing.T) {
	if err := UpdateGuildHelpRequestCountByID(int64(9070260011401216), 9); err != nil {
		t.Fatal(err)
	}
}

func Test_BatchGetLatestReqeustTimeByGuildIDs(t *testing.T) {
	_, err := BatchGetLatestReqeustTimeByGuildIDs([]int64{9068658676465664, 9194785835319296})
	if err != nil {
		t.Fatal(err)
	}
}
