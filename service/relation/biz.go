package main

import (
	"context"
	pbmessage "douyin/grpc/message"
	pbrelation "douyin/grpc/relation"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"sync"

	"go.uber.org/zap"
)

type RelationService struct {
	pbrelation.UnimplementedRelationServer
}

func (s *RelationService) Follow(ctx context.Context, req *pbrelation.FollowRequest) (*pbrelation.FollowResponse, error) {
	if req.UserID == req.ToUserID {
		// 这里要去考虑 到底是直接返回error 还是
		return &pbrelation.FollowResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.CantNotFollowSelf,
		}, nil
	}
	// 需要先看看有没有关注 不能重复关注
	isFollow, err := IsFollow(req.UserID, req.ToUserID)
	if err != nil { // 缓存不存在去查库
		zap.L().Sugar().Warn(constant.CacheMiss)
		isFollow, err = IsFollowed(req.UserID, req.ToUserID)
		if err != nil {
			zap.L().Sugar().Error(err)
			return &pbrelation.FollowResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.DatabaseError,
			}, err
		}
		// 异步更新缓存
		go func() {
			followIDs, err := SelectFollowingByUserID(req.UserID)
			if err != nil {
				zap.L().Sugar().Error(err)
				return
			}
			err = SetFollowUserIDSet(req.UserID, followIDs)
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
	err = Follow(req.UserID, req.ToUserID, 1)
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

func (s *RelationService) Unfollow(ctx context.Context, req *pbrelation.FollowRequest) (*pbrelation.FollowResponse, error) {
	if req.UserID == req.ToUserID {
		return &pbrelation.FollowResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.CantNotUnFollowSelf,
		}, nil
	}
	isFollow, err := IsFollow(req.UserID, req.ToUserID)
	if err != nil { // 缓存不存在去查库
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		isFollow, err = IsFollowed(req.UserID, req.ToUserID)
		if err != nil {
			zap.L().Sugar().Error(err)
			return &pbrelation.FollowResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.DatabaseError,
			}, nil
		}
		// 异步更新缓存
		go func() {
			followIDs, err := SelectFollowingByUserID(req.UserID)
			if err != nil {
				zap.L().Sugar().Error(err)
				return
			}
			err = SetFollowUserIDSet(req.UserID, followIDs)
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
	SendFollowMessage(req.UserID, req.ToUserID, -1)
	zap.L().Debug("Follow: publish a message")
	return &pbrelation.FollowResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.UnFollowSuccess,
	}, nil
}

func (s *RelationService) FollowList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.ListResponse, error) {
	// 先使用布隆过滤器判断userID 存不存在
	// FIXME 这里需要考虑的是 布隆过滤器在哪里用
	// 网关层 还是具体到服务里 ？？ 感觉网关会好点
	// 或者用redis实现一个 又或者每次RPC调用？？
	// if !UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
	// 	err := fmt.Errorf(constant.BloomFilterRejected)
	// 	zap.L().Error(err.Error())
	// 	return nil, err
	// }
	following, err := GetFollowUserIDSet(req.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		following, err = SelectFollowingByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
		// 将缓存写入
		err = SetFollowUserIDSet(req.UserID, following)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
	usersInfo := make([]*pbuser.UserInfo, len(following))
	wg := &sync.WaitGroup{}
	wg.Add(len(following))
	for i := range following {
		go func(i int) {
			defer wg.Done()
			u, err := userClient.Info(ctx, &pbuser.InfoRequest{
				LoginUserID: req.LoginUserID,
				UserID:      following[i],
			})
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			usersInfo[i] = u.GetUser()
		}(i)
	}
	wg.Wait()
	return &pbrelation.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FollowListSuccess,
		UserList:   usersInfo,
	}, nil
}

func (s *RelationService) FollowerList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.ListResponse, error) {
	// if !UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
	// 	err := fmt.Errorf(constant.BloomFilterRejected)
	// 	zap.L().Error(err.Error())
	// 	return nil, err
	// }
	follower, err := GetFollowerUserIDSet(req.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		follower, err = SelectFollowerByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = SetFollowerUserIDSet(req.UserID, follower)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	usersInfo := make([]*pbuser.UserInfo, len(follower))
	wg := &sync.WaitGroup{}
	wg.Add(len(follower))
	for i := range follower {
		go func(i int) {
			defer wg.Done()
			u, err := userClient.Info(ctx, &pbuser.InfoRequest{
				LoginUserID: req.LoginUserID,
				UserID:      follower[i],
			})
			if err != nil {
				zap.L().Error(err.Error())
			}
			usersInfo[i] = u.GetUser()
		}(i)
	}
	wg.Wait()
	return &pbrelation.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FollowerListSuccess,
		UserList:   usersInfo,
	}, nil
}

// 好友列表 互相关注即为好友
// 客户端需要知道最近的一条消息
func (s *RelationService) FriendList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.FriendsResponse, error) {
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
	following, err := GetFollowUserIDSet(req.UserID)
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		following, err = SelectFollowingByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = SetFollowUserIDSet(req.UserID, following)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	follower, err := GetFollowerUserIDSet(req.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		follower, err = SelectFollowerByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = SetFollowerUserIDSet(req.UserID, follower)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	// 判断用户的粉丝列表有没有关注
	followerMap := make(map[uint64]struct{}, len(follower))
	for _, item := range follower {
		followerMap[item] = struct{}{}
	}
	friends := make([]uint64, 0)
	for _, ff := range following {
		if _, ok := followerMap[ff]; ok {
			friends = append(friends, ff)
		}
	}
	friends = append(friends, constant.ChatGPTID)
	friendsInfo := make([]*pbrelation.FriendInfo, len(friends))
	wg := &sync.WaitGroup{}
	wg.Add(len(friends))
	for i := range friends {
		go func(i int) {
			defer wg.Done()
			friendInfo, err := userClient.Info(ctx, &pbuser.InfoRequest{
				LoginUserID: req.LoginUserID,
				UserID:      friends[i],
			})
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			msg, err := messageClient.GetFirstMessage(ctx, &pbmessage.GetFirstRequest{
				UserID:   req.UserID,
				ToUserID: friends[i],
			})
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			friend := &pbrelation.FriendInfo{
				Id:              friendInfo.User.Id,
				Name:            friendInfo.User.Name,
				Avatar:          friendInfo.User.Avatar,
				BackgroundImage: friendInfo.User.BackgroundImage,
				Signature:       friendInfo.User.Signature,
				IsFollow:        friendInfo.User.IsFollow,
				FollowCount:     friendInfo.User.FollowCount,
				FollowerCount:   friendInfo.User.FollowerCount,
				TotalFavorited:  friendInfo.User.TotalFavorited,
				WorkCount:       friendInfo.User.WorkCount,
				FavoriteCount:   friendInfo.User.FavoriteCount,
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

func (s *RelationService) IsFollow(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.IsFollowResponse, error) {
	isFollow, err := IsFollow(req.LoginUserID, req.UserID)
	// 缓存未命中 查询数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		isFollow, err = IsFollowed(req.LoginUserID, req.UserID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return &pbrelation.IsFollowResponse{
				Result: false,
			}, nil
		}
		go func() {
			// 关注列表
			followUserIDSet, err := SelectFollowingByUserID(req.LoginUserID)
			if err != nil {
				zap.L().Error(constant.DatabaseError, zap.Error(err))
				return
			}
			err = SetFollowUserIDSet(req.LoginUserID, followUserIDSet)
			if err != nil {
				zap.L().Sugar().Error(err)
			}
		}()
	}
	return &pbrelation.IsFollowResponse{
		Result: isFollow,
	}, nil
}

func (s *RelationService) IsFriend(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.IsFriendResponse, error) {
	isFollow, err := IsFollow(req.LoginUserID, req.UserID)
	// 缓存未命中 查询数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		isFollow, err = IsFollowed(req.LoginUserID, req.UserID)
		if err != nil {
			zap.L().Error(constant.DatabaseError, zap.Error(err))
			return &pbrelation.IsFriendResponse{
				Result: false,
			}, nil
		}
		go func() {
			// 关注列表
			followUserIDSet, err := SelectFollowingByUserID(req.LoginUserID)
			if err != nil {
				zap.L().Error(constant.DatabaseError, zap.Error(err))
				return
			}
			err = SetFollowUserIDSet(req.LoginUserID, followUserIDSet)
			if err != nil {
				zap.L().Sugar().Error(err)
			}
		}()
	}
	if isFollow {
		isFollowed, err := IsFollow(req.UserID, req.LoginUserID)
		// 缓存未命中 查询数据库
		if err != nil {
			zap.L().Warn(constant.CacheMiss, zap.Error(err))
			isFollowed, err = IsFollowed(req.UserID, req.LoginUserID)
			if err != nil {
				zap.L().Error(constant.DatabaseError, zap.Error(err))
				return &pbrelation.IsFriendResponse{
					Result: false,
				}, nil
			}
			go func() {
				// 关注列表
				followUserIDSet, err := SelectFollowingByUserID(req.UserID)
				if err != nil {
					zap.L().Error(constant.DatabaseError, zap.Error(err))
					return
				}
				err = SetFollowUserIDSet(req.UserID, followUserIDSet)
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
