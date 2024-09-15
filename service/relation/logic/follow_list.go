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

func FollowList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.ListResponse, error) {
	// 先使用布隆过滤器判断userID 存不存在
	// FIXME 这里需要考虑的是 布隆过滤器在哪里用
	// 网关层 还是具体到服务里 ？？ 感觉网关会好点
	// 或者用redis实现一个 又或者每次RPC调用？？
	// if !UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
	// 	err := fmt.Errorf(constant.BloomFilterRejected)
	// 	zap.L().Error(err.Error())
	// 	return nil, err
	// }
	following, err := storage.GetFollowUserIDSet(req.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		following, err = storage.SelectFollowingByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
		// 将缓存写入
		err = storage.SetFollowUserIDSet(req.UserID, following)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
	usersInfo := make([]*pbuser.UserInfo, 0, len(following))
	userMap, err := rpc.UserClient.List(ctx, &pbuser.ListReq{
		LoginUserID: req.LoginUserID,
		UserID:      following,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	for _, u := range following {
		usersInfo = append(usersInfo, userMap.User[u])
	}
	return &pbrelation.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FollowListSuccess,
		UserList:   usersInfo,
	}, nil
}
