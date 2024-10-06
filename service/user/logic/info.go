package logic

import (
	"context"
	"errors"
	"strconv"

	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/service/user/storage"
	"diktok/storage/cache"
	"diktok/storage/database/model"

	"go.uber.org/zap"
)

func Info(ctx context.Context, req *pbuser.InfoRequest) (*pbuser.InfoResponse, error) {
	// 使用布隆过滤器判断用户ID是否存在
	if !cache.UserIDBloomFilter.TestString(strconv.FormatInt(req.UserID, 10)) {
		err := errors.New(constant.BloomFilterRejected)
		zap.L().Sugar().Error(err)
		return nil, err
	}
	// 去redis里查询用户信息 这是热点数据 redis缓存确实快了很多
	user, err := storage.GetUserInfo(req.UserID)
	// 缓存未命中再去查数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		user, err = storage.SelectUserByID(req.UserID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return nil, err
		}
		// 设置缓存
		go func() {
			err = storage.SetUserInfo(user)
			if err != nil {
				zap.L().Error(constant.SetCacheError, zap.Error(err))
			}
		}()
	}
	// 判断是否是关注用户
	var isFollow bool
	// 用户未登录
	if req.LoginUserID == 0 {
		isFollow = false
	} else if req.LoginUserID == req.UserID { // 自己查自己 当然是关注了的
		isFollow = true
	} else {
		isFollowRes, err := rpc.RelationClient.IsFollow(ctx, &pbrelation.ListRequest{
			UserID:      req.UserID,
			LoginUserID: req.LoginUserID,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		isFollow = isFollowRes.Result
	}
	return &pbuser.InfoResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.LoadSuccess,
		User:       userResponse(user, isFollow),
	}, nil
}

func userResponse(user *model.User, isFollowed bool) *pbuser.UserInfo {
	return &pbuser.UserInfo{
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Id:              user.ID,
		IsFollow:        isFollowed,
		Name:            user.Username,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
}
