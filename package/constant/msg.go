package constant

// 消息
const (
	// 状态码 TODO 要去除依赖
	Success = 0
	Failed  = -1

	// zap.L().Warn(constant.CacheMiss, zap.Error(err)) 统一使用
	CacheMiss           = "缓存未命中"
	SetCacheError       = "设置缓存失败"
	BloomFilterRejected = "布隆过滤器拦截"
	DatabaseError       = "数据库操作失败"
	BadParaRequest      = "参数错误"
	LoadSuccess         = "加载成功"

	// 用户
	SecretFormatError = "密码格式错误"
	SecretFormatEasy  = "密码太简单"
	UserDepulicate    = "用户名已存在"
	UsernameFormatErr = "用户名格式错误"
	FrequentLogin     = "登录次数过多 5分钟后再试"
	UserNoExist       = "用户不存在"
	SecretError       = "用户密码错误"
	RegisterSuccess   = "用户注册成功"
	LoginSuccess      = "用户登录成功"
	TooLongSignature  = "签名长度不符合要求"
	UpdateSuccess     = "更新用户信息成功"

	// 视频
	PublishListSuccess = "发布视频列表获取成功"
	SearchSuccess      = "搜索成功"
	FeedSuccess        = "视频列表获取成功"
	NoMoreVideos       = "视频见底了"
	UploadVideoSuccess = "上传视频成功"
	GetTokenSuccess    = "获取上传凭证成功"

	// relation
	FollowSuccese       = "关注成功"
	FollowListSuccess   = "关注列表加载成功"
	FollowerListSuccess = "粉丝列表加载成功"
	CantNotFollowSelf   = "不能关注自己"
	FollowRepeated      = "不能重复关注"
	FollowError         = "关注失败"
	UnFollowError       = "取关失败"
	UnFollowSuccess     = "取关成功"
	CantNotUnFollowSelf = "不能取关自己"
	UnFollowNotFollowed = "不能取关未关注的人"
	DefaultMessage      = "开始对话吧"
	FriendListError     = "无法查看别人的好友列表"
	FriendListSuccess   = "好友列表加载成功"

	// 评论
	CommentSuccess       = "评论成功"
	DeleteCommentSuccess = "删除评论成功"
	LoadCommentsSuccess  = "加载评论列表成功"

	// 点赞
	FavoriteSuccess     = "点赞成功"
	UnFavoriteSuccess   = "取消点赞成功"
	FavoriteListSuccess = "喜欢视频列表获取成功"

	// message
	SendSuccess = "发送成功"
	ListSuccess = "消息列表加载成功"
	SendToSelf  = "不能给自己发送消息"
	SendEmpty   = "消息内容为空"
	IsNotFriend = "对方不是你的好友"
)
