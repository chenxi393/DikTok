package constant

import (
	"time"
)

// 一些常量
const (
	DoAction           = "1"
	UndoAction         = "2"
	UserID             = "userID"
	DebugMode          = "debug"
	EasySecret         = "123456"
	SnoyFlakeStartTime = 1698775594477
	MaxVideoNumber     = 30
	DefaultCover       = "default_cover.png"
	ServiceName        = "diktok"

	// topic字段 前后端都是写死的目前
	TopicDefualt = "现在短视频非常的流行热门"
	TopicSport   = "体育"
	TopicGame    = "游戏"
	TopicMusic   = "音乐"

	// chatgpt
	ChatGPTAvatar = "2022chatgpt.png"
	ChatGPTName   = "ChatGPT"
	ChatGPTID     = 1

	// token过期时间
	TokenExpiration = 48 * time.Hour
)
