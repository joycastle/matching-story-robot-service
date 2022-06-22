package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Guild 公会表
type Guild struct {
	ID          int64  `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"index:,class:FULLTEXT,option:WITH PARSER ngram;size:50;not null" json:"name"`
	Description string `gorm:"size:200" json:"description"`
	Country     string `json:"country"`

	LevelLimit int32  `json:"level_limit"`
	Badge      string `json:"badge"`

	PresidentID int64

	Score          int64  `json:"score"`
	Star           int64  `gorm:"not null;default:0" json:"star"`
	StarActivityID int64  `gorm:"column:star_activity_id" json:"star_activity_id"`
	Help           int32  `gorm:"default:0" json:"help"`
	HelpWeekID     int32  `gorm:"default:0" json:"help_week_id"` // 公会帮助每周清零，故该ID记录help所属的周数。
	FrameIcon      string `gorm:"default:0" json:"frame_icon"`
	Robot          bool
	MemberNum      int32 `gorm:"column:member_num"`
	IsOpen         int8

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name of the Guild model
func (g *Guild) TableName() string {
	return "guild"
}

// JsonFormat returns JSON format after Marshal
func (g *Guild) JsonFormat() string {
	jsonStr, _ := json.Marshal(g)
	return string(jsonStr)
}
