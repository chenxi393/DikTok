package response

type RelationListResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户信息列表
	UserList []User `json:"user_list"`
}

type FriendResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户信息列表
	UserList []FriendUser `json:"user_list"`
}
type FriendUser struct {
	User
	Message string `json:"message"`
	MsgType int  `json:"msgType"`
}
