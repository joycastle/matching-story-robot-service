package service

import (
	"testing"
)

func Test_GetGuildResponeUidsByRequestID(t *testing.T) {
	request_id := int64(9181379573055488)
	list, err := GetGuildResponeUidsByRequestID(request_id)
	if err != nil {
		t.Fatal(err, list)
	}

	if list, err := GetGuildResoneRobotUserByRequestID(int64(73802617859342336)); err != nil {
		t.Fatal(err, list)
	}
}

func Test_InsertOne(t *testing.T) {
	requestID := int64(9079666878971904)
	rspUserID := int64(35469000060)
	reqUserID := int64(35469000060)

	if r, err := AddGuildHelpResone(requestID, rspUserID, reqUserID); err != nil {
		t.Fatal(err, r)
	} else if r.ID <= 0 {
		t.Fatal("")
	}
}
