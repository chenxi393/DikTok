package response

const (
	Success             = 0
	Failed              = -1
	RegisterSuccess     = "用户注册成功"
	LoginSucess         = "用户登录成功"
	BadParaRequest      = "参数错误，请求失败"
	WrongToken          = "token不匹配 请重新登录"
	FeedSuccess         = "视频列表获取成功"
	NoMoreVideos        = "视频见底了"
	FileFormatError     = "文件格式错误"
	UploadVideoSuccess  = "上传视频成功"
	PubulishListSuccess = "发布视频列表获取成功"
	FavoriteListSuccess = "喜欢视频列表获取成功"
	FollowListSuccess   = "关注列表加载成功"
	FollowerListSuccess = "粉丝列表加载成功"
	ActionSuccess       = "操作成功"
	FriendListError     = "无法查看别人的好友列表"
	FriendListSuccess   = "好友列表加载成功"
	SendSuccess         = "发送成功"
)

type CommonResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
