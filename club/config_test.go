package club

import (
	"testing"
	"time"
)

func Test_AAA(t *testing.T) {
	if getRobotMaxLimitNum() <= 0 || getNormalUserNum() <= 0 {
		t.Fatal("")
	}

	ok := false
	for i := 0; i < 200; i++ {
		num := getGenerateRobotNumByRand()
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
		num := getLikeNumByRand()
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

}
