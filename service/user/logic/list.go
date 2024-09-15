package logic

import (
	"context"

	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/service/user/storage"
	"diktok/storage/database/model"

	"go.uber.org/zap"
)

func List(ctx context.Context, req *pbuser.ListReq) (*pbuser.ListResp, error) {
	userMap := make(map[int64]*model.User)
	for _, u := range req.GetUserID() {
		// 去redis里查询用户信息 这是热点数据 redis缓存确实快了很多
		user, err := storage.GetUserInfo(u)
		// 缓存未命中再去查数据库
		if err != nil {
			zap.L().Warn(constant.CacheMiss, zap.Error(err))
			user, err = storage.SelectUserByID(u)
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
		userMap[u] = user
	}

	respMap := make(map[int64]*pbuser.UserInfo, len(userMap))
	for _, u := range userMap {
		// 判断是否是关注用户
		var isFollow bool
		// 用户未登录
		if req.LoginUserID == 0 {
			isFollow = false
		} else if req.LoginUserID == u.ID { // 自己查自己 当然是关注了的
			isFollow = true
		} else {
			isFollowRes, err := rpc.RelationClient.IsFollow(ctx, &pbrelation.ListRequest{
				UserID:      u.ID,
				LoginUserID: req.LoginUserID,
			})
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
			isFollow = isFollowRes.Result
		}
		respMap[u.ID] = userResponse(u, isFollow)
	}
	return &pbuser.ListResp{
		User: respMap,
	}, nil
}
