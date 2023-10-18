package service

import (
	"douyin/database"
	"douyin/response"
	"fmt"

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

func (service *RelationService) RelationAction(userID uint64) error {
	if userID == service.ToUserID {
		err := fmt.Errorf("不能关注和取消关注自己")
		zap.L().Info(err.Error())
		return err
	}
	// 关注
	err := fmt.Errorf("ActionType 错误")
	if service.ActionType == "1" {
		err = database.Follow(userID, service.ToUserID, 1)
	} else if service.ActionType == "2" {
		err = database.Follow(userID, service.ToUserID, -1)
	}
	return err
}

func (service *RelationService) RelationFollowList(userID uint64) (*response.RelationListResponse, error) {
	// 先拿出所有的关注用户ID 再去用户表拿出所有信息 再判断有没有关注
	following, err := database.SelectFollowingByUserID(service.UserID) //是service的UserID
	if err != nil {
		return nil, err
	}
	users, err := database.SelectUserListByIDs(following)
	if err != nil {
		return nil, err
	}
	// 拿用户的关注列表
	loginUserFollowing, err := database.SelectFollowingByUserID(userID)
	if err != nil {
		return nil, err
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
		StatusMsg:  "关注列表加载成功",
		UserList:   useersResponse,
	}
	return res, nil
}

func (service *RelationService) RelationFollowerList(userID uint64) (*response.RelationListResponse, error) {
	// 先拿出所有的粉丝用户ID 再去用户表拿出所有信息 再判断有没有关注
	follower, err := database.SelectFollowerByUserID(service.UserID)
	if err != nil {
		return nil, err
	}
	// 拿粉丝列表的用户信息
	users, err := database.SelectUserListByIDs(follower)
	if err != nil {
		return nil, err
	}
	// 拿用户的关注列表
	following, err := database.SelectFollowingByUserID(userID)
	if err != nil {
		return nil, err
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
		StatusMsg:  "粉丝列表加载成功",
		UserList:   useersResponse,
	}
	return res, nil
}

// 好友列表 互相关注即为好友
// TODO 这个api文档好像不对 客户端似乎需要知道最近的一条消息
func (service *RelationService) RelationFriendList() (*response.FriendResponse, error) {
	//应该是可以用自连接的
	//或者先去关注表找自己已经关注的人 然后去关注表 找关注自己的人 （可以用in）
	//拿到一组ID 就是好友 然后再批量拿出信息
	following, err := database.SelectFollowingByUserID(service.UserID)
	if err != nil {
		return nil, err
	}
	follower, err := database.SelectFollowerByUserID(service.UserID)
	if err != nil {
		return nil, err
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
	// 拿好友的信息
	friendsInfo, err := database.SelectUserListByIDs(friends)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// TODO 所以还要去消息列表拿最近的一条消息 接口文档有误
	usersResponse := make([]response.FriendUser, 0, len(friends))
	for i := range friendsInfo {
		// fix 这里循环查库了 记得规避
		msg, err := database.GetMessageNewest(service.UserID, friendsInfo[i].ID)
		if err != nil || msg == "" {
			msg = "快来开启和好友的第一次对话吧！！！"
		}
		uu := response.FriendUser{
			User:    *response.UserInfo(&friendsInfo[i], true),
			Message: msg, // 最近的一条消息 客户端测试了是有的
			MsgType: 0,
		}
		usersResponse = append(usersResponse, uu)
	}

	res := &response.FriendResponse{
		StatusCode: response.Success,
		StatusMsg:  "好友列表加载成功",
		UserList:   usersResponse,
	}
	return res, nil
}
