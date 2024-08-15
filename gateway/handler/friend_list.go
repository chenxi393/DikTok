package handler

import (
	"diktok/gateway/response"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type friendListRequest struct {
	// 用户id List使用 查看这个用户的关注列表，粉丝列表，好友列表
	UserID int64 `query:"user_id"`
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
