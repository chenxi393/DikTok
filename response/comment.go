package response

type CommentActionResponse struct {
	// 评论成功返回评论内容，不需要重新拉取整个列表
	Comment *Comment `json:"comment"`
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
}

type CommentListResponse struct {
	// 评论列表
	CommentList []Comment `json:"comment_list"`
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
}

type Comment struct {
	// 评论内容
	Content string `json:"content"`
	// 评论发布日期，格式 yyyy-mm-dd
	CreateDate string `json:"create_date"`
	// 评论id
	ID uint64 `json:"id"`
	// 评论用户信息
	User User `json:"user"`
}

// func BuildComment(comment *model.Comment, user *model.User, isFollow bool) *Comment {
// 	return &Comment{
// 		Content:    comment.Content,
// 		// 这个评论的时间客户端哈好像可以到毫秒2006-01-02 15:04:05.999
// 		// 但是感觉每必要 日期就够了
// 		CreateDate: comment.CreatedTime.Format("2006-01-02 15:04"),
// 		ID:         comment.ID,
// 		User:       *UserInfo(user, isFollow),
// 	}
// }
