package handler

import (
	"diktok/gateway/response"
	pbfavorite "diktok/grpc/favorite"
	pbvideo "diktok/grpc/video"
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
	if len(Favoriteresp.VideoList) <= 0 {
		return c.JSON(response.BuildVideoList(nil))
	}
	packResp, err := rpc.VideoClient.Pack(c.UserContext(), &pbvideo.PackReq{
		LoginUserId: loginUserID,
		VideoId:     Favoriteresp.GetVideoList(),
	})
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}
	return c.JSON(response.BuildVideoList(packResp.GetVideoList()))
}
