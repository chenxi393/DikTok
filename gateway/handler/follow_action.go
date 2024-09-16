package handler

import (
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
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
		return c.JSON(constant.InvalidParams)
	}
	userID := c.Locals(constant.UserID).(int64)
	var res *pbrelation.FollowResponse
	if req.ActionType == constant.DoAction {
		res, err = rpc.RelationClient.Follow(c.UserContext(), &pbrelation.FollowRequest{
			UserID:   userID,
			ToUserID: req.ToUserID,
		})
	} else if req.ActionType == constant.UndoAction {
		res, err = rpc.RelationClient.Unfollow(c.UserContext(), &pbrelation.FollowRequest{
			UserID:   userID,
			ToUserID: req.ToUserID,
		})
	} else {
		err = constant.InvalidParams
	}
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
