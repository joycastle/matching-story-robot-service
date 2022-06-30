package service

import "testing"

func TestUpdateRobotByRobotID(t *testing.T) {
	if err := UpdateRobotByRobotID("club_130714000009", "group_id", 999, "act_num", "1"); err != nil {
		t.Fatal(err)
	}
}
