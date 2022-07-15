package lib

import (
	"errors"
	"fmt"
	"testing"
)

func TestLogStructuredJson(t *testing.T) {
	log := NewLogStructed()
	log.Failed().Step(1).Err(fmt.Errorf("hhhhhhhh")).Set("a", "b", 1, 10, "d", 199, "f", 192.22)
	fmt.Println(log.String())
	fmt.Println(log.Success().Step(0).String())
	fmt.Println(log.Success().Step(0).Set("error", errors.New("XXXXXXX")).String())
	fmt.Println(log.Success().Step(0).Set("error", fmt.Errorf("XXXXXXX")).String())
}
