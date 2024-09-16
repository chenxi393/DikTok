package handler

import (
	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func FavoriteList(c *fiber.Ctx) error {
	var req likeListRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.InvalidParams)
	}
	loginUserID := c.Locals(constant.UserID).(int64)
	Favoriteresp, err := rpc.FavoriteClient.List(c.UserContext(), &pbfavorite.ListRequest{
		UserID: req.UserID,
	})
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}
	resp, err := BuildVideosInfo(c.Context(), Favoriteresp.GetVideoList(), nil, loginUserID)
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}
	return c.JSON(resp)
}
