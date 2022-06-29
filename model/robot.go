package model

import (
	"encoding/json"
)

// Robot 机器人表结构
type Robot struct {
	ID int64 `json:"id" gorm:"column:id;autoIncrement;uniqueIndex"`
	// 机器人引用ID
	RobotID    string `json:"robot_id" gorm:"column:robot_id;primaryKey"`
	ConfID     int32  `json:"conf_id" gorm:"column:conf_id;index"`
	ActivityID string `json:"activity_id" gorm:"column:activity_id;default:0;index"`
	GroupID    int64  `json:"group_id" gorm:"column:group_id;default:0;index"`
	IndexNum   int32  `json:"index_num" gorm:"column:index_num;default:0"`

	RankNum int64 `json:"rank_num" gorm:"column:rank_num;default:0"`

	Name     string `json:"name" gorm:"column:name;not null"`
	HeadIcon string `json:"head_icon" gorm:"column:head_icon;not null"`

	CreateTime int64 `json:"create_time" gorm:"column:create_time;default:0"`
	// 活动结束时间
	EndTime int64 `json:"end_time" gorm:"column:end_time;default:0"`

	ActNum int32 `json:"act_num" gorm:"column:act_num;default:0"`

	// 下次行动时间
	NextActTime int64 `json:"next_act_time" gorm:"column:next_act_time;default:0"`
	// 下次行动加的积分
	NextPassNum int64 `json:"next_pass_num" gorm:"column:next_pass_num;default:0"`
	// 机器人最高分数
	RobotMax int64 `json:"robot_max" gorm:"column:robot_max;default:0"`
	// 当前阶段
	Period int32 `json:"period" gorm:"column:period;default:0"`
	// 当前阶段分数
	PeriodScore int64 `json:"period_score" gorm:"column:period_score;default:0"`
	// 当前阶段最高分
	PeriodMaxScore  int64  `json:"period_max_score" gorm:"column:period_max_score;default:0"`
	ActivityGroupID string `json:"activity_group_id" gorm:"column:activity_group_id;default:0;index"`
}

// TableName returns the table name of the User model
func (p *Robot) TableName() string {
	return "robot_table"
}

// JsonFormat returns JSON format after Marshal
func (p *Robot) JsonFormat() string {
	jsonStr, _ := json.Marshal(p)
	return string(jsonStr)
}
