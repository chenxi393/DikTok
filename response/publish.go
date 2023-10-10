package response

type PublishResponse struct {
	StatusCode int  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
