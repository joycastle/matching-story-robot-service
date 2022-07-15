package create

type Job struct {
	GuildID    int64 `json:"guild_id"`
	ActionTime int64 `json:"action_time"`
	IsInit     bool  `json:"is_init"`
}
