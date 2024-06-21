package handler

import (
	"errors"

	"diktok/gateway/response"
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
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
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
		err = errors.New(constant.BadParaRequest)
	}
	if err != nil {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func FavoriteList(c *fiber.Ctx) error {
	var req likeListRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.VideoListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
			VideoList:  nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	resp, err := FavoriteClient.List(c.UserContext(), &pbfavorite.ListRequest{
		UserID:      req.UserID,
		LoginUserID: userID,
	})
	if err != nil {
		res := response.VideoListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
			VideoList:  nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
