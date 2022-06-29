package action

import (
	"errors"
	"fmt"
	"time"

	"github.com/joycastle/matching-story-robot-service/club/config"
	"github.com/joycastle/matching-story-robot-service/service"
)

var (
	weedDaysConfig map[time.Weekday]int = make(map[time.Weekday]int, 7)
)

func init() {
	weedDaysConfig[time.Sunday] = 7
	weedDaysConfig[time.Monday] = 1
	weedDaysConfig[time.Tuesday] = 2
	weedDaysConfig[time.Wednesday] = 3
	weedDaysConfig[time.Thursday] = 4
	weedDaysConfig[time.Friday] = 5
	weedDaysConfig[time.Saturday] = 6
}

func robotActionBeforeCheck(job *Job) error {

	robotConfig, err := service.GetRobotForGuild(job.UserID)
	if err != nil {
		return err
	}

	actionID := robotConfig.GroupID
	activeDaysMap := config.GetRobotActiveDaysByActionID(int(actionID))

	if len(activeDaysMap) == 0 {
		return errors.New(fmt.Sprintf("active day config not found active_id:%d", actionID))
	}

	if _, ok := activeDaysMap[-1]; ok && len(activeDaysMap) == 1 {
		return errors.New("active day config is [-1], not take effect")
	}

	todayWeek := time.Now().Weekday()
	todayWeekInt := weedDaysConfig[todayWeek]
	if _, ok := activeDaysMap[todayWeekInt]; !ok {
		return errors.New(fmt.Sprintf("today:week:%d is not a active day, %v, active_id:%d", todayWeekInt, activeDaysMap, actionID))
	}

	//判断是否被踢出工会
	if u, err := service.GetGuildInfoByIDAndUid(job.GuildID, job.UserID); err != nil {
		return err
	} else if u.GuildID <= 0 {
		return errors.New(fmt.Sprintf("already kick out the guild"))
	}

	return nil
}
