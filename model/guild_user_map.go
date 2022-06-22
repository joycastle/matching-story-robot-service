package model

import (
	"encoding/json"
	"time"
)

type GuildUserMap struct {
	UserID  int64 `gorm:"primary_key" json:"user_id"`
	GuildID int64 `gorm:"index:guild_id_idx"`

	CreatedAt time.Time
}

// TableName returns the table name of the GuildUserMap model
func (g *GuildUserMap) TableName() string {
	return "guild_user_map"
}

// JsonFormat returns JSON format after Marshal
func (g *GuildUserMap) JsonFormat() string {
	jsonStr, _ := json.Marshal(g)
	return string(jsonStr)
}
