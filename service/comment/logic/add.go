package logic

import (
	"context"
	pbcomment "diktok/grpc/comment"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/comment/storage"
	"diktok/storage/database/model"

	"time"

	"go.uber.org/zap"
)

func Add(ctx context.Context, req *pbcomment.AddRequest) (*pbcomment.CommentResponse, error) {
	// TODO 增加敏感词过滤 可以异步实现 comment表多一列屏蔽信息
	id, err := util.GetSonyFlakeID()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	msg := &model.Comment{
		ID:          int64(id),
		VideoID:     req.VideoID,
		UserID:      req.UserID,
		ParentID:    req.ParentID,
		Content:     req.Content,
		CreatedTime: time.Now(),
		ToUserID:    req.ToUserID,
	}
	err = storage.CreateComment(msg)
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
		Id:         msg.ID,
		User:       userResponse.GetUser()[req.UserID],
		Content:    msg.Content,
		CreateDate: msg.CreatedTime.Format("2006-01-02 15:04"),
	}
	return &pbcomment.CommentResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.CommentSuccess,
		Comment:    commentResponse,
	}, nil
}
