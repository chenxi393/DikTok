package handler

import (
	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	FavoriteClient pbfavorite.FavoriteClient
)

type likeRequest struct {
	// 1-点赞，2-取消点赞
	ActionType string `query:"action_type"`
	// 视频id
	VideoID int64 `query:"video_id"`
}

type likeListRequest struct {
	UserID int64 `query:"user_id"`
}

func FavoriteVideoAction(c *fiber.Ctx) error {
	var req likeRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.InvalidParams)
	}
	userID := c.Locals(constant.UserID).(int64)
	var resp *pbfavorite.LikeResponse
	if req.ActionType == constant.DoAction {
		resp, err = FavoriteClient.Like(c.UserContext(), &pbfavorite.LikeRequest{
			UserID:  userID,
			VideoID: req.VideoID,
		})
	} else if req.ActionType == constant.UndoAction {
		resp, err = FavoriteClient.Unlike(c.UserContext(), &pbfavorite.LikeRequest{
			UserID:  userID,
			VideoID: req.VideoID,
		})
	} else {
		err = constant.InvalidParams
	}
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}
	return c.JSON(resp)
}
