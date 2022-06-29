package service

import (
	"fmt"
	"testing"
)

func TestSendChatMessageV1(t *testing.T) {
	if ret, err := SendChatMessageRPC("44531B0BCB34A58BBB1D9F92CA5330B7", 130290000000, 125323777000079360, "你好"); err != nil {
		t.Fatal(err, ret)
	} else {
		fmt.Println(ret)
	}
}

func TestSendRequestHelpRPC(t *testing.T) {
	if ret, err := SendRequestHelpRPC("44531B0BCB34A58BBB1D9F92CA5330B7", 130290000000, 125323777000079360, 125361495293820928); err != nil {
		t.Fatal(err, ret)
	}
}
