package model

import "encoding/json"

// User 用户表结构
type User struct {
	ID              uint   `json:"id" gorm:"column:id;type:int unsigned auto_increment;not null;uniqueIndex:id_unique_idx"`
	UserID          int64  `json:"user_id" gorm:"primaryKey;column:user_id;type:bigint"`
	AccountID       string `json:"account_id" gorm:"column:account_id;type:varchar(64);not null;index:account_id"`
	UserName        string `json:"user_name" gorm:"column:user_name;type:varchar(64);not null"`
	UserHeadIcon    string `json:"user_head_icon" gorm:"column:user_head_icon;type:varchar(512)"`
	Device          string `json:"device" gorm:"column:device;type:varchar(64)"`
	DeviceType      int    `json:"device_type" gorm:"column:device_type;type:int;not null;default:0"`
	Status          int    `json:"status" gorm:"index;column:status;type:int;not null;default:0"`
	LoginTime       int64  `json:"login_time" gorm:"column:login_time;type:bigint;not null;default:0"`
	LoginDays       []byte `json:"login_days" gorm:"column:login_days;type:blob"`
	Channel         int32  `json:"channel" gorm:"column:channel;type:tinyint;not null;default:0"`
	Language        int32  `json:"language" gorm:"column:language;type:tinyint;not null;default:0"`
	TotalPay        int64  `json:"total_pay" gorm:"column:total_pay;type:bigint;not null;default:0"`
	Mails           []byte `json:"mails" gorm:"column:mails;type:blob"`
	MailIDs         []byte `json:"mail_ids" gorm:"column:mail_ids;type:blob"`
	Flag            int    `json:"flag" gorm:"column:flag;type:int;not null;default:0"`
	UserData        []byte `json:"user_data" gorm:"column:user_data;type:mediumblob"`
	UserFriendData  []byte `json:"user_friend_data" gorm:"column:user_friend_data;type:mediumblob"`
	UserLevel       int    `json:"user_level" gorm:"column:user_level;type:int;not null;default:1;index:user_level_idx"`
	UserLikeData    []byte `json:"user_like_data" gorm:"column:user_like_data;type:mediumblob"`
	UserUnlockID    int    `json:"user_unlock_id" gorm:"column:user_unlock_id;type:int;not null;default:0"`
	UserBuyData     []byte `json:"user_buy_data" gorm:"column:user_buy_data;type:blob"`
	UserCountryData string `json:"user_country_data" gorm:"column:user_country_data;type:varchar(64)"`
	CheatUser       int32  `json:"cheat_user" gorm:"column:cheat_user;type:tinyint;not null;default:0"`
	ClientVersion   int    `json:"client_version" gorm:"column:client_version;type:int;not null;default:0"`
	CreateTime      int64  `json:"create_time" gorm:"column:create_time;type:bigint;not null;default:0"`
	UpdateTime      int64  `json:"update_time" gorm:"column:update_time;type:bigint;not null;default:0"`
	DeleteTime      int64  `json:"delete_time" gorm:"column:delete_time;type:bigint;not null;default:0"`
	UserStar        int    `json:"user_star" gorm:"column:user_star;type:int;not null;default:0"`
	StarActivityID  int64  `gorm:"column:star_activity_id" json:"star_activity_id"`
	UserHelp        int    `json:"user_help" gorm:"column:user_help;type:int;not null;default:0"`
	TestUser        int    `json:"test_user" gorm:"column:test_user;type:int;not null;default:0"`

	// 需要在存档时从redis缓存存储数据
	UserLikeCount       uint   `json:"user_like_count" gorm:"column:user_like_count;type:int unsigned;not null;default:0"`
	ApplyFriendData     []byte `json:"apply_friend_data" gorm:"column:apply_friend_data;type:mediumblob"`
	LikedFriendData     []byte `json:"liked_friend_data" gorm:"column:liked_friend_data;type:mediumblob"`
	RecommendFriendData []byte `json:"recommend_friend_data" gorm:"column:recommend_friend_data;type:mediumblob"`

	UserType string `json:"user_type" gorm:"column:user_type;type:varchar(16);default:normal"`
}

// TableName returns the table name of the User model
func (u *User) TableName() string {
	return "user_table"
}

// JsonFormat returns JSON format after Marshal
func (u *User) JsonFormat() string {
	jsonStr, _ := json.Marshal(u)
	return string(jsonStr)
}
