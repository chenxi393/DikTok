package constant

import (
	"time"
)

// 一些常量
const (
	TokenTimeOut    = 12 * time.Hour
	TokenMaxRefresh = 3 * time.Hour

	DoAction           = "1"
	UndoAction         = "2"
	UserID             = "userID"
	SnoyFlakeStartTime = 1698775594477
)

// redis 的key
const (
	// 登录次数的key 默认10分钟 登录5次 + user_name
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

// 一些redis过期时间
const (
	MaxLoginTime    = 5
	MaxloginInernal = 5 * time.Minute
	Expiration      = 300 * time.Second
)

// 消息
const (
	// zap.L().Warn(constant.CacheMiss, err) 统一使用
	CacheMiss           = "缓存未命中"
	DatabaseError       = "数据库操作失败"
	BloomFilterRejected = "布隆过滤器拦截"
	BadParaRequest      = "参数错误"
	SetCacheError       = "设置缓存失败"
	EmptyKey            = "键不存在"

	// 评论
	CommentSuccess       = "评论成功"
	DeleteCommentSuccess = "删除评论成功"
	LoadCommentsSuccess  = "加载评论列表成功"

	// 点赞
	FavoriteSuccess   = "点赞成功"
	UnFavoriteSuccess = "取消点赞成功"
)
