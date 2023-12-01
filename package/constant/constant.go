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

	// topic字段 前后端都是写死的目前
	TopicDefualt = "现在短视频非常的流行热门"
	TopicSport   = "体育"
	TopicGame    = "游戏"
	TopicMusic   = "音乐"
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
	TokenTimeOut = 12 * time.Hour
)

// 消息
const (
	// zap.L().Warn(constant.CacheMiss, zap.Error(err)) 统一使用
	CacheMiss           = "缓存未命中"
	SetCacheError       = "设置缓存失败"
	EmptyKey            = "键不存在"
	BloomFilterRejected = "布隆过滤器拦截"
	DatabaseError       = "数据库操作失败"
	BadParaRequest      = "参数错误"

	// 评论
	CommentSuccess       = "评论成功"
	DeleteCommentSuccess = "删除评论成功"
	LoadCommentsSuccess  = "加载评论列表成功"

	// 点赞
	FavoriteSuccess   = "点赞成功"
	UnFavoriteSuccess = "取消点赞成功"

	// 用户
	SecretFormatError = "密码格式错误"
	SecretFormatEasy  = "密码太简单"
	UserDepulicate    = "用户名已被注册"
	FrequentLogin     = "登录次数过多 5分钟后再试"
	UserNoExist       = "用户不存在"
	SecretError       = "用户密码错误"

	// 视频
	VideoServerBug = "视频缺少作者 服务端bug"

	// relation
	CantNotFollowSelf   = "不能关注自己"
	FollowError         = "关注失败"
	UnFollowError       = "取关失败"
	CantNotUnFollowSelf = "不能取关自己"
	UnFollowError1      = "不能取关未关注的人"
	DefaultMessage      = "开始对话吧"
)
