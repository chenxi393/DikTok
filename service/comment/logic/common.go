package logic

type CommentExtra struct {
	ToCommentID int64  `json:"to_comment_id,omitempty"` // 回复的评论ID
	ImageURI    string `json:"image_uri,omitempty"`
}
