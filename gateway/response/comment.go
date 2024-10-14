package response

import (
	pbcomment "diktok/grpc/comment"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/util"
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
	Content         string `json:"content"`     // 评论内容
	CreateDate      string `json:"create_date"` // 评论发布日期，格式 yyyy-mm-dd
	CommentID       int64  `json:"id"`
	ParentID        int64  `json:"parent_id,omitempty"`         // 根评论id
	ImageURL        string `json:"image_url,omitempty"`         // 图片
	ToCommentID     int64  `json:"to_comment_id,omitempty"`     // 回复某条评论ID
	SubCommentCount int64  `json:"sub_comment_count,omitempty"` // 子评论数量
	User            *User  `json:"user"`                        // 评论用户信息
}

func BuildCommentList(comList *pbcomment.ListResponse, userList *pbuser.ListResp, countResp *pbcomment.CountResp) *CommentListResponse {
	res := &CommentListResponse{
		CommentList: make([]*Comment, 0, len(comList.GetCommentList())),
		HasMore:     comList.GetHasMore(),
		Total:       comList.GetTotal(),
		StatusCode:  constant.Success,
		StatusMsg:   constant.LoadCommentsSuccess,
	}

	UserMap := BuildUserMap(userList)
	for _, v := range comList.GetCommentList() {
		if v != nil {
			res.CommentList = append(res.CommentList, BuildComment(v, UserMap, countResp.GetCountMap()))
		}
	}
	return res
}

func BuildComment(v *pbcomment.CommentData, userMp map[int64]*User, countMap map[int64]int64) *Comment {
	return &Comment{
		User:       userMp[v.UserID],
		Content:    v.Content,
		CommentID:  v.CommentID,
		CreateDate: time.Unix(v.CreateAt, 0).Format(time.DateTime),

		ParentID:        v.ParentID,
		ImageURL:        util.Uri2Url(v.ImageURI),
		ToCommentID:     v.ToCommentID,
		SubCommentCount: countMap[v.CommentID],
	}
}
