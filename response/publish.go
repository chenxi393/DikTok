package response

type PublishListResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户发布的视频列表
	VideoList []Video `json:"video_list"`
}
