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
	MP4Suffix          = ".mp4"
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
)

// redis 的key
const (
	// 登录次数的key 默认5分钟 登录5次 + user_name
	LoginCounterPrefix = "login_counter:"

	// + user_id
	UserInfoPrefix = "user_info:"
	// + user_id
	UserInfoCountPrefix = "user_info_count:"
	// + video_id
	VideoInfoPrefix = "video_info:"
	// + video_id
	VideoInfoCountPrefix = "video_info_count:"

	// user 哈希hset 的键
	FollowCountField    = "follow_count:"
	FollowerCountField  = "follower_count:"
	TotalFavoritedField = "total_favorited_count:"
	WorkCountField      = "work_count:"
	FavoriteCountField  = "favorite_count:"

	// video 哈希hset 的键
	FavoritedCountField = "favorited_count:"
	CommentCountField   = "comment_count"

	// + user_id
	FollowIDPrefix = "follow_id:"
	// + user_id
	FollowerIDPrefix = "follower_id:"
	// + user_id
	FavoriteIDPrefix = "favorite_id:"
	// + user_id
	PublishIDPrefix = "publish_id:"
	// + video_id
	CommentPrefix = "comment:"
)

// 过期时间
const (
	// redis过期时间
	MaxLoginTime    = 5
	MaxloginInernal = 5 * time.Minute
	Expiration      = 400 * time.Second

	// 分布式锁
	LockTime  = 200 // 200ms 以毫秒为单位
	RetryTime = 20

	// token过期时间
	TokenExpiration = 48 * time.Hour
)

// grpc
var (
	UserService     = "diktok/user"
	VideoService    = "diktok/video"
	RalationService = "diktok/relation"
	CommentService  = "diktok/comment"
	MessageService  = "diktok/message"
	FavoriteService = "diktok/favorite"
	VideoAddr       = "video:8010" //docker 内使用video:8010 本地使用127.0.0.1 否则会阻塞
	UserAddr        = "user:8020"
	RelationAddr    = "relation:8030"
	MessageAddr     = "message:8040"
	FavoriteAddr    = "favorite:8050"
	CommentAddr     = "comment:8060"
)
