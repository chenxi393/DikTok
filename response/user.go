package response

type UserRegisterOrLogin struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户鉴权token
	Token string `json:"token"`
	// 用户id
	UserID uint64 `json:"user_id"`
}
