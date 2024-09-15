package logic

import (
	"context"

	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/service/relation/storage"

	"go.uber.org/zap"
)

func FollowerList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.ListResponse, error) {
	// if !UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
	// 	err := fmt.Errorf(constant.BloomFilterRejected)
	// 	zap.L().Error(err.Error())
	// 	return nil, err
	// }
	follower, err := storage.GetFollowerUserIDSet(req.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		follower, err = storage.SelectFollowerByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = storage.SetFollowerUserIDSet(req.UserID, follower)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	usersInfo := make([]*pbuser.UserInfo, 0, len(follower))
	userMap, err := rpc.UserClient.List(ctx, &pbuser.ListReq{
		LoginUserID: req.LoginUserID,
		UserID:      follower,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	for _, u := range follower {
		usersInfo = append(usersInfo, userMap.User[u])
	}
	return &pbrelation.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FollowerListSuccess,
		UserList:   usersInfo,
	}, nil
}
