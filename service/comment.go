package service

import (
	"douyin/database"
	"douyin/model"
	"douyin/response"
	"fmt"

	"go.uber.org/zap"
)

type CommentService struct {
	// 1-发布评论，2-删除评论
	ActionType string `query:"action_type"`
	// 要删除的评论id，在action_type=2的时候使用
	CommentID *string `query:"comment_id,omitempty"`
	// 用户填写的评论内容，在action_type=1的时候使用
	CommentText *string `query:"comment_text,omitempty"`
	// 用户鉴权token
	Token string `query:"token"`
	// 视频id
	VideoID uint64 `query:"video_id"`
}

func (service *CommentService) CommentAction(userID uint64) (*response.CommentActionResponse, error) {
	err := fmt.Errorf("参数错误")
	var comment *model.Comment
	if service.ActionType == "1" && service.CommentText != nil {
		// 发布评论
		// TODO 增加敏感词过滤 我觉得可以异步
		// TODO 基于雪花算法 生成评论ID 为什么 主键自增不可以吗
		// TODO 好像POST请求 一般都使用消息队列
		comment, err = database.CommentAdd(service.CommentText, service.VideoID, userID)
	} else if service.ActionType == "2" && service.CommentID != nil {
		comment, err = database.CommentDelete(service.CommentID, service.VideoID, userID)
	}
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 查找评论的用户信息
	user, err := database.SelectUserByID(userID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 这里的 isFollow 直接返回 false ，因为评论人自己当然不能关注自己
	// TODO 没懂为什么是false 看看客户端把
	var msg string
	if service.ActionType == "1" {
		msg = "评论成功"
	} else {
		msg = "删除成功"
	}
	return &response.CommentActionResponse{
		StatusCode: response.Success,
		StatusMsg:  msg,
		Comment:    response.BuildComment(comment, user, false),
	}, nil
}

func (service *CommentService) CommentList(userID uint64) (*response.CommentListResponse, error) {
	// TODO 使用布隆过滤器判断视频ID是否存在
	// Redis啥的
	// 先拿到这个视频的所有评论
	comment, err := database.GetCommentsByVideoID(service.VideoID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 再去拿每一个评论的作者信息
	userIDs := make([]uint64, 0, len(comment))
	for _, cc := range comment {
		userIDs = append(userIDs, cc.UserID)
	}
	// 范围查询
	users, err := database.SelectUserListByIDs(userIDs)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 判断这些用户是否被关注
	followingIDs, err := database.SelectFollowingByUserID(userID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	followingMap := make(map[uint64]struct{}, len(followingIDs))
	for _, ff := range followingIDs {
		followingMap[ff] = struct{}{}
	}
	// 再把user放到map里 避免双重循环去查找F
	usersMap := make(map[uint64]*model.User, len(users))
	for i, user := range users {
		usersMap[user.ID] = &users[i]
	}
	// 构造返回值
	commentResponse := make([]response.Comment, 0, len(comment))
	for _, cc := range comment {
		_, isFollowed := followingMap[cc.UserID]
		res := response.BuildComment(cc, usersMap[cc.UserID], isFollowed)
		commentResponse = append(commentResponse, *res)
	}
	return &response.CommentListResponse{
		StatusCode:  response.Success,
		StatusMsg:   "加载评论列表成功",
		CommentList: commentResponse,
	}, nil
}
