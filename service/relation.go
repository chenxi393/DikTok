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
		zap.L().Error(err.Error())
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
	useersResponse := make([]response.User, 0, len(users))
	for i := range users {
		uu := response.UserInfo(&users[i], true)
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
