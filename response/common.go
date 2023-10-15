package response

var (
	Success         = 0
	Failed          = -1
	RegisterSuccess = "用户注册成功"
	LoginSucess     = "用户登录成功"
	BadParaRequest  = "参数错误，请求失败"
	WrongToken      = "token不匹配 请重新登录"
	FeedSuccess     = "视频列表获取成功"
)

type CommonResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
