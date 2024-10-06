package response

import (
	pbcomment "diktok/grpc/comment"
	"time"
)

type CommentActionResponse struct {
	// 评论成功返回评论内容，不需要重新拉取整个列表
	// 按道理这里不应该返回 让前端自己渲染
	Comment *Comment `json:"comment"`
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
}

type CommentListResponse struct {
	// 评论列表
	CommentList []*Comment `json:"comment_list"`
	HasMore     bool       `json:"has_more"`
	Total       int64      `json:"total"`
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
	CommentID int64 `json:"id"`
	// 评论用户信息
	User *User `json:"user"`
	// TODO 记得增加信息
}

func BuildComment(v *pbcomment.CommentData, userMp map[int64]*User) *Comment {
	return &Comment{
		User:       userMp[v.UserID],
		Content:    v.Content,
		CommentID:  v.CommentID,
		CreateDate: time.Unix(v.CreateAt, 0).Format(time.DateTime),
	}
}
