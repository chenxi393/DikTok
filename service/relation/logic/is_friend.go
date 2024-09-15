package logic

import (
	"context"

	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/service/relation/storage"

	"go.uber.org/zap"
)

func IsFriend(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.IsFriendResponse, error) {
	isFollow, err := storage.IsFollow(req.LoginUserID, req.UserID)
	// 缓存未命中 查询数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		isFollow, err = storage.IsFollowed(req.LoginUserID, req.UserID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return &pbrelation.IsFriendResponse{
				Result: false,
			}, nil
		}
		go func() {
			// 关注列表
			followUserIDSet, err := storage.SelectFollowingByUserID(req.LoginUserID)
			if err != nil {
				zap.L().Error(constant.DatabaseError, zap.Error(err))
				return
			}
			err = storage.SetFollowUserIDSet(req.LoginUserID, followUserIDSet)
			if err != nil {
				zap.L().Sugar().Error(err)
			}
		}()
	}
	if isFollow {
		isFollowed, err := storage.IsFollow(req.UserID, req.LoginUserID)
		// 缓存未命中 查询数据库
		if err != nil {
			zap.L().Warn(constant.CacheMiss, zap.Error(err))
			isFollowed, err = storage.IsFollowed(req.UserID, req.LoginUserID)
			if err != nil {
				zap.L().Error(constant.DatabaseError, zap.Error(err))
				return &pbrelation.IsFriendResponse{
					Result: false,
				}, nil
			}
			go func() {
				// 关注列表
				followUserIDSet, err := storage.SelectFollowingByUserID(req.UserID)
				if err != nil {
					zap.L().Error(constant.DatabaseError, zap.Error(err))
					return
				}
				err = storage.SetFollowUserIDSet(req.UserID, followUserIDSet)
				if err != nil {
					zap.L().Sugar().Error(err)
				}
			}()
		}
		return &pbrelation.IsFriendResponse{
			Result: isFollowed,
		}, nil
	}
	return &pbrelation.IsFriendResponse{
		Result: isFollow,
	}, nil
}
