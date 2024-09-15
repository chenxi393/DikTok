package logic

import (
	"context"

	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/service/relation/storage"

	"go.uber.org/zap"
)

func Follow(ctx context.Context, req *pbrelation.FollowRequest) (*pbrelation.FollowResponse, error) {
	if req.UserID == req.ToUserID {
		// 这里要去考虑 到底是直接返回error 还是
		return &pbrelation.FollowResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.CantNotFollowSelf,
		}, nil
	}
	// 需要先看看有没有关注 不能重复关注
	isFollow, err := storage.IsFollow(req.UserID, req.ToUserID)
	if err != nil { // 缓存不存在去查库
		zap.L().Sugar().Warn(constant.CacheMiss)
		isFollow, err = storage.IsFollowed(req.UserID, req.ToUserID)
		if err != nil {
			zap.L().Sugar().Error(err)
			return &pbrelation.FollowResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.DatabaseError,
			}, err
		}
		// 异步更新缓存
		go func() {
			followIDs, err := storage.SelectFollowingByUserID(req.UserID)
			if err != nil {
				zap.L().Sugar().Error(err)
				return
			}
			err = storage.SetFollowUserIDSet(req.UserID, followIDs)
			if err != nil {
				zap.L().Sugar().Error(err)
			}
		}()
	}
	if isFollow {
		return &pbrelation.FollowResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FollowRepeated,
		}, nil
	}
	// 放入消息队列 异步 返回成功
	// 这里先写入redis 再写入数据库
	err = storage.Follow(req.UserID, req.ToUserID, 1)
	if err != nil {
		zap.L().Sugar().Error(err)
		return &pbrelation.FollowResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.DatabaseError,
		}, nil
	}
	zap.L().Debug("Follow: publish a message")
	return &pbrelation.FollowResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FollowSuccese,
	}, nil
}
