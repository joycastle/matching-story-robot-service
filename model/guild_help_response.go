package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GuildHelpResponse struct {
	ID int64 `gorm:"primary_key" json:"id"`

	HelpID int64 `gorm:"index:help_id_idx"`

	ResponderID int64 `gorm:"index:responder_id_idx"`
	RequesterID int64 `gorm:"index:requester_id_idx"`

	Ack bool `gorm:"index:ack_idx" json:"ack"` // 是否被消费

	Time int64 `json:"time"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name of the GuildHelpResponse model
func (g *GuildHelpResponse) TableName() string {
	return "guild_help_response"
}

// JsonFormat returns JSON format after Marshal
func (g *GuildHelpResponse) JsonFormat() string {
	jsonStr, _ := json.Marshal(g)
	return string(jsonStr)
}
