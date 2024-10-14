package logic

import (
	"context"

	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/service/relation/storage"

	"go.uber.org/zap"
)

func Unfollow(ctx context.Context, req *pbrelation.FollowRequest) (*pbrelation.FollowResponse, error) {
	if req.UserID == req.ToUserID {
		return &pbrelation.FollowResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.CantNotUnFollowSelf,
		}, nil
	}
	isFollow, err := storage.IsFollow(req.UserID, req.ToUserID)
	if err != nil { // 缓存不存在去查库
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		isFollow, err = storage.IsFollowed(req.UserID, req.ToUserID)
		if err != nil {
			zap.L().Sugar().Error(err)
			return &pbrelation.FollowResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.DatabaseError,
			}, nil
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
	if !isFollow {
		return &pbrelation.FollowResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.UnFollowNotFollowed,
		}, nil
	}
	storage.SendFollowMessage(req.UserID, req.ToUserID, -1)
	zap.L().Sugar().Debugf("Follow: publish a message")
	return &pbrelation.FollowResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.UnFollowSuccess,
	}, nil
}
