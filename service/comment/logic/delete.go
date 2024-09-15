package logic

import (
	"context"

	pbcomment "diktok/grpc/comment"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/service/comment/storage"
)

func Delete(ctx context.Context, req *pbcomment.DeleteRequest) (*pbcomment.CommentResponse, error) {
	// 我们认为删除评论不是高频动作 故不使用消息队列
	// database里会删缓存 并且校验是不是自己发的 实际上不校验也行
	// 注意还需要在database里减少视频的评论数
	msg, err := storage.DeleteComment(req.CommentID, req.VideoID, req.UserID)
	if err != nil {
		return nil, err
	}
	// 查找评论的用户信息
	userResponse, err := rpc.UserClient.List(ctx, &pbuser.ListReq{
		UserID:      []int64{req.UserID},
		LoginUserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	commentResponse := &pbcomment.CommentData{
		Id:      msg.ID,
		User:    userResponse.GetUser()[req.UserID],
		Content: msg.Content,
		// 这个评论的时间客户端哈好像可以到毫秒2006-01-02 15:04:05.999
		// 但是感觉每必要 日期就够了
		CreateDate: msg.CreatedTime.Format("2006-01-02 15:04"),
	}
	return &pbcomment.CommentResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.DeleteCommentSuccess,
		Comment:    commentResponse,
	}, nil
}
