package handler

import (
	"context"
	"douyin/gateway/auth"
	"douyin/gateway/response"
	pbrelation "douyin/grpc/relation"
	"douyin/package/constant"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	RelationClient pbrelation.RelationClient
)

type followRequest struct {
	// 1-关注，2-取消关注
	ActionType string `query:"action_type"`
	// 对方用户id
	ToUserID uint64 `query:"to_user_id"`
}

type followListRequest struct {
	Token string `query:"token"`
	// 用户id List使用 查看这个用户的关注列表，粉丝列表，好友列表
	UserID uint64 `query:"user_id"`
}

type friendListRequest struct {
	// 用户id List使用 查看这个用户的关注列表，粉丝列表，好友列表
	UserID uint64 `query:"user_id"`
}

func RelationAction(c *fiber.Ctx) error {
	var req followRequest
	err := c.QueryParser(&req)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	var res *pbrelation.FollowResponse
	if req.ActionType == constant.DoAction {
		res, err = RelationClient.Follow(context.Background(), &pbrelation.FollowRequest{
			UserID:   userID,
			ToUserID: req.ToUserID,
		})
	} else if req.ActionType == constant.UndoAction {
		res, err = RelationClient.Unfollow(context.Background(), &pbrelation.FollowRequest{
			UserID:   userID,
			ToUserID: req.ToUserID,
		})
	} else {
		err = errors.New(constant.BadParaRequest)
	}
	if err != nil {
		// 这里由于rpc会传递具体的错误信息
		// 可以考虑不用
		res := &response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

func FollowList(c *fiber.Ctx) error {
	var req followListRequest
	err := c.QueryParser(&req)
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var userID uint64
	if req.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(req.Token)
		if err != nil {
			res := response.UserRegisterOrLogin{
				StatusCode: constant.Failed,
				StatusMsg:  constant.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		userID = claims.UserID
	}
	resp, err := RelationClient.FollowList(context.Background(), &pbrelation.ListRequest{
		LoginUserID: userID,
		UserID:      req.UserID,
	})
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func FollowerList(c *fiber.Ctx) error {
	var req followListRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var userID uint64
	if req.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(req.Token)
		if err != nil {
			res := response.UserRegisterOrLogin{
				StatusCode: constant.Failed,
				StatusMsg:  constant.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		userID = claims.UserID
	}
	resp, err := RelationClient.FollowerList(context.Background(), &pbrelation.ListRequest{
		LoginUserID: userID,
		UserID:      req.UserID,
	})
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func FriendList(c *fiber.Ctx) error {
	var req friendListRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	if userID != req.UserID {
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FriendListError,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := RelationClient.FriendList(context.Background(), &pbrelation.ListRequest{
		LoginUserID: userID,
		UserID:      req.UserID,
	})
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
