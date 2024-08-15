package handler

import (
	"errors"

	"diktok/gateway/response"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
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
