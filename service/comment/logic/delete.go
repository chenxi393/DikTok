package logic

import (
	"context"

	pbbase "diktok/grpc/base"
	pbcomment "diktok/grpc/comment"
	"diktok/package/constant"
	"diktok/service/comment/storage"
)

func Delete(ctx context.Context, req *pbcomment.DeleteRequest) (*pbbase.BaseResp, error) {
	// 有子评论耗时 500ms
	// 没有耗时 300ms
	// 怎么优化
	err := storage.DeleteCommentContent(ctx, req.CommentID)
	if err != nil {
		return nil, err
	}
	return &pbbase.BaseResp{
		StatusCode: constant.Success,
		StatusMsg:  constant.DeleteCommentSuccess,
	}, nil
}
