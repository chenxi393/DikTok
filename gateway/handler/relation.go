package handler

import (
	"errors"

	"diktok/gateway/response"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"

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
	ToUserID int64 `query:"to_user_id"`
}

type followListRequest struct {
	// 用户id List使用 查看这个用户的关注列表，粉丝列表，好友列表
	UserID int64 `query:"user_id"`
}

type friendListRequest struct {
	// 用户id List使用 查看这个用户的关注列表，粉丝列表，好友列表
	UserID int64 `query:"user_id"`
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
	userID := c.Locals(constant.UserID).(int64)
	var res *pbrelation.FollowResponse
	if req.ActionType == constant.DoAction {
		res, err = RelationClient.Follow(c.UserContext(), &pbrelation.FollowRequest{
			UserID:   userID,
			ToUserID: req.ToUserID,
		})
	} else if req.ActionType == constant.UndoAction {
		res, err = RelationClient.Unfollow(c.UserContext(), &pbrelation.FollowRequest{
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
	userID := c.Locals(constant.UserID).(int64)
	resp, err := RelationClient.FollowList(c.UserContext(), &pbrelation.ListRequest{
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
	userID := c.Locals(constant.UserID).(int64)
	resp, err := RelationClient.FollowerList(c.UserContext(), &pbrelation.ListRequest{
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
	userID := c.Locals(constant.UserID).(int64)
	if userID != req.UserID {
		res := response.RelationListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FriendListError,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := RelationClient.FriendList(c.UserContext(), &pbrelation.ListRequest{
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
