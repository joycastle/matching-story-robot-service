package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GuildHelpRequest struct {
	ID int64 `gorm:"primary_key" json:"id"`

	Total    int16
	Count    int16
	Resource int8 // 请求帮助的资源，体力/金币/碎片

	GuildID int64 `gorm:"index:guild_id_idx"`

	RequesterID int64

	Done bool  `gorm:"default:false" json:"done"` // 是否已经完成
	Time int64 `json:"time"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name of the GuildHelpRequest model
func (g *GuildHelpRequest) TableName() string {
	return "guild_help_request"
}

// JsonFormat returns JSON format after Marshal
func (g *GuildHelpRequest) JsonFormat() string {
	jsonStr, _ := json.Marshal(g)
	return string(jsonStr)
}
