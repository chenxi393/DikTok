package service

import (
	"douyin/database"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/llm"
	"douyin/package/mq"
	"douyin/response"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

type RelationService struct {
	// 1-关注，2-取消关注
	ActionType string `query:"action_type"`
	// 对方用户id
	ToUserID uint64 `query:"to_user_id"`
	// 用户鉴权token
	Token string `query:"token"`
	// 用户id List使用 查看这个用户的关注列表，粉丝列表，好友列表
	UserID uint64 `query:"user_id"`
}

func (service *RelationService) FollowAction(userID uint64) error {
	if userID == service.ToUserID {
		return errors.New(constant.CantNotFollowSelf)
	}
	// 需要先看看有没有关注 不能重复关注
	isFollow, err := cache.IsFollow(userID, service.ToUserID)
	if err != nil { // 缓存不存在去查库
		zap.L().Sugar().Warn(constant.CacheMiss)
		isFollow, err = database.IsFollowed(userID, service.ToUserID)
		if err != nil {
			zap.L().Sugar().Error(err)
			return err
		}
		// 异步更新缓存
		go func() {
			//followIDs, err := database.SelectFollowingByUserID(userID)
			if err != nil {
				zap.L().Sugar().Error(err)
				return
			}
			//err = cache.SetFollowUserIDSet(userID, followIDs)
			if err != nil {
				zap.L().Sugar().Error(err)
			}
		}()
	}
	if isFollow {
		return errors.New(constant.FollowError)
	}
	// 放入消息队列 异步 返回成功
	// 这里先写入redis 再写入数据库
	mq.SendFollowMessage(userID, service.ToUserID, 1)
	zap.L().Debug("Follow: publish a message")
	return nil
}

func (service *RelationService) UnFollowAction(userID uint64) error {
	if userID == service.ToUserID {
		err := errors.New(constant.CantNotUnFollowSelf)
		zap.L().Sugar().Error(err)
		return err
	}
	isFollow, err := cache.IsFollow(userID, service.ToUserID)
	if err != nil { // 缓存不存在去查库
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		isFollow, err = database.IsFollowed(userID, service.ToUserID)
		if err != nil {
			zap.L().Sugar().Error(err)
			return err
		}
		// 异步更新缓存
		go func() {
			//followIDs, err := database.SelectFollowingByUserID(userID)
			if err != nil {
				zap.L().Sugar().Error(err)
				return
			}
			//err = cache.SetFollowUserIDSet(userID, followIDs)
			if err != nil {
				zap.L().Sugar().Error(err)
			}
		}()
	}
	if !isFollow {
		err := errors.New(constant.UnFollowError1)
		zap.L().Sugar().Error(err)
		return err
	}
	mq.SendFollowMessage(userID, service.ToUserID, -1)
	zap.L().Debug("Follow: publish a message")
	return nil
}

func (service *RelationService) RelationFollowList(userID uint64) (*response.RelationListResponse, error) {
	// 先使用布隆过滤器判断userID 存不存在
	if !cache.UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
		err := fmt.Errorf(constant.BloomFilterRejected)
		zap.L().Error(err.Error())
		return nil, err
	}
	following, err := cache.GetFollowUserIDSet(service.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		following, err = database.SelectFollowingByUserID(service.UserID)
		if err != nil {
			return nil, err
		}
		// 将缓存写入
		err = cache.SetFollowUserIDSet(service.UserID, following)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
	users, err := database.SelectUserListByIDs(following)
	if err != nil {
		return nil, err
	}
	// 拿用户的关注列表
	var loginUserFollowing []uint64
	if userID != 0 {
		loginUserFollowing, err = cache.GetFollowUserIDSet(userID)
		if err != nil {
			zap.L().Warn(constant.CacheMiss, zap.Error(err))
			loginUserFollowing, err = database.SelectFollowingByUserID(userID)
			if err != nil {
				return nil, err
			}
			go func() {
				err := cache.SetFollowUserIDSet(userID, loginUserFollowing)
				if err != nil {
					zap.L().Error(err.Error())
				}
			}()
		}
	}
	// 判断用户的关注列表登录用户有没有关注
	loginUserFollowingMap := make(map[uint64]struct{}, len(following))
	for _, item := range loginUserFollowing {
		loginUserFollowingMap[item] = struct{}{}
	}
	loginUserFollowingMap[userID] = struct{}{}
	useersResponse := make([]response.User, 0, len(users))
	for i, user := range users {
		_, ok := loginUserFollowingMap[user.ID]
		uu := response.UserInfo(&users[i], ok)
		useersResponse = append(useersResponse, *uu)
	}
	res := &response.RelationListResponse{
		StatusCode: response.Success,
		StatusMsg:  response.FollowListSuccess,
		UserList:   useersResponse,
	}
	return res, nil
}

func (service *RelationService) RelationFollowerList(userID uint64) (*response.RelationListResponse, error) {
	if !cache.UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
		err := fmt.Errorf(constant.BloomFilterRejected)
		zap.L().Error(err.Error())
		return nil, err
	}
	follower, err := cache.GetFollowerUserIDSet(service.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		follower, err = database.SelectFollowerByUserID(service.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = cache.SetFollowerUserIDSet(service.UserID, follower)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	// 拿粉丝列表的用户信息
	users, err := database.SelectUserListByIDs(follower)
	if err != nil {
		return nil, err
	}
	// 拿用户的关注列表
	var following []uint64
	if userID != 0 {
		following, err = cache.GetFollowUserIDSet(userID)
		if err != nil {
			zap.L().Sugar().Warn(constant.CacheMiss)
			following, err = database.SelectFollowingByUserID(userID)
			if err != nil {
				return nil, err
			}
			go func() {
				// 将缓存写入
				err = cache.SetFollowUserIDSet(userID, following)
				if err != nil {
					zap.L().Error(err.Error())
				}
			}()
		}
	}
	// 判断用户的粉丝列表有没有关注
	followingMap := make(map[uint64]struct{}, len(following))
	for _, item := range following {
		followingMap[item] = struct{}{}
	}
	// 自己也在自己的关注列表里
	followingMap[userID] = struct{}{}
	useersResponse := make([]response.User, 0, len(users))
	for i, user := range users {
		_, ok := followingMap[user.ID]
		uu := response.UserInfo(&users[i], ok)
		useersResponse = append(useersResponse, *uu)
	}
	res := &response.RelationListResponse{
		StatusCode: response.Success,
		StatusMsg:  response.FollowerListSuccess,
		UserList:   useersResponse,
	}
	return res, nil
}

// 好友列表 互相关注即为好友
// 客户端需要知道最近的一条消息
func (service *RelationService) RelationFriendList() (*response.FriendResponse, error) {
	if !cache.UserIDBloomFilter.TestString(strconv.FormatUint(service.UserID, 10)) {
		err := fmt.Errorf(constant.BloomFilterRejected)
		zap.L().Error(err.Error())
		return nil, err
	}
	//应该是可以用自连接的
	//或者先去关注表找自己已经关注的人 然后去关注表 找关注自己的人 （可以用in）
	//拿到一组ID 就是好友 然后再批量拿出信息
	// 拿用户的关注列表
	following, err := cache.GetFollowUserIDSet(service.UserID)
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		following, err = database.SelectFollowingByUserID(service.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = cache.SetFollowUserIDSet(service.UserID, following)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	follower, err := cache.GetFollowerUserIDSet(service.UserID)
	// 缓存未命中 查数据库
	if err != nil {
		zap.L().Sugar().Warn(constant.CacheMiss)
		follower, err = database.SelectFollowerByUserID(service.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			// 将缓存写入
			err = cache.SetFollowerUserIDSet(service.UserID, follower)
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
	friends = append(friends, llm.ChatGPTID)
	// 拿好友的信息
	friendsInfo, err := database.SelectUserListByIDs(friends)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	usersResponse := make([]response.FriendUser, 0, len(friends))
	for i := range friendsInfo {
		// FIXME这里循环查库了 记得规避
		msg, err := database.GetMessageNewest(service.UserID, friendsInfo[i].ID)
		msgt := 0
		if err != nil || msg.Content == "" {
			msg.Content = constant.DefaultMessage
		} else {
			if msg.FromUserID == service.UserID {
				msgt = 1
			}
		}
		uu := response.FriendUser{
			User:    *response.UserInfo(&friendsInfo[i], true),
			Message: msg.Content, // 最近的一条消息 客户端测试了是有的
			MsgType: msgt,
		}
		usersResponse = append(usersResponse, uu)
	}
	res := &response.FriendResponse{
		StatusCode: response.Success,
		StatusMsg:  response.FriendListSuccess,
		UserList:   usersResponse,
	}
	return res, nil
}
