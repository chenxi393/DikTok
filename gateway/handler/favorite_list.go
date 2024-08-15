package handler

import (
	"diktok/gateway/response"
	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func FavoriteList(c *fiber.Ctx) error {
	var req likeListRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(response.BuildStdResp(constant.Failed, constant.BadParaRequest, nil))
	}
	loginUserID := c.Locals(constant.UserID).(int64)
	Favoriteresp, err := FavoriteClient.List(c.UserContext(), &pbfavorite.ListRequest{
		UserID: req.UserID,
	})
	if err != nil {
		return c.JSON(response.BuildStdResp(constant.Failed, err.Error(), nil))
	}
	resp, err := BuildVideosInfo(c.Context(), Favoriteresp.GetVideoList(), nil, loginUserID)
	if err != nil {
		return c.JSON(response.BuildStdResp(constant.Failed, constant.InternalException, nil))
	}
	return c.JSON(resp)
}
