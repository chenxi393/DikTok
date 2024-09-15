package logic

import (
	"context"
	"sync"

	pbmessage "diktok/grpc/message"
	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/service/relation/storage"

	"go.uber.org/zap"
)

// 好友列表 互相关注即为好友
// 客户端需要知道最近的一条消息
func FriendList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.FriendsResponse, error) {
	// 布隆过滤器这块看看怎么用
	// if !UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
	// 	err := fmt.Errorf(constant.BloomFilterRejected)
	// 	zap.L().Error(err.Error())
	// 	return nil, err
	// }
	//应该是可以用自连接的
	//或者先去关注表找自己已经关注的人 然后去关注表 找关注自己的人 （可以用in）
	//拿到一组ID 就是好友 然后再批量拿出信息
	// 拿用户的关注列表
	following, err := storage.GetFollowUserIDSet(req.UserID)
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		following, err = storage.SelectFollowingByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = storage.SetFollowUserIDSet(req.UserID, following)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
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
	// 判断用户的粉丝列表有没有关注
	followerMap := make(map[int64]struct{}, len(follower))
	for _, item := range follower {
		followerMap[item] = struct{}{}
	}
	friends := make([]int64, 0)
	for _, ff := range following {
		if _, ok := followerMap[ff]; ok {
			friends = append(friends, ff)
		}
	}
	friends = append(friends, constant.ChatGPTID)
	friendsInfo := make([]*pbrelation.FriendInfo, len(friends))
	userMap, err := rpc.UserClient.List(ctx, &pbuser.ListReq{
		LoginUserID: req.LoginUserID,
		UserID:      friends,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(friends))
	for i := range friends {
		go func(i int) {
			defer wg.Done()
			friendInfo := userMap.User[friends[i]]
			msg, err := rpc.MessageClient.GetFirstMessage(ctx, &pbmessage.GetFirstRequest{
				UserID:   req.UserID,
				ToUserID: friends[i],
			})
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			friend := &pbrelation.FriendInfo{
				Id:              friendInfo.Id,
				Name:            friendInfo.Name,
				Avatar:          friendInfo.Avatar,
				BackgroundImage: friendInfo.BackgroundImage,
				Signature:       friendInfo.Signature,
				IsFollow:        friendInfo.IsFollow,
				FollowCount:     friendInfo.FollowCount,
				FollowerCount:   friendInfo.FollowerCount,
				TotalFavorited:  friendInfo.TotalFavorited,
				WorkCount:       friendInfo.WorkCount,
				FavoriteCount:   friendInfo.FavoriteCount,
				Message:         msg.Message, // 最近的一条消息 客户端测试了是有的
				MsgType:         msg.MsgType,
			}
			friendsInfo[i] = friend
		}(i)
	}
	wg.Wait()
	return &pbrelation.FriendsResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FriendListSuccess,
		UserList:   friendsInfo,
	}, nil
}
