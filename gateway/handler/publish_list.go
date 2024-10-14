package handler

import (
	"diktok/gateway/response"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type listRequest struct {
	UserID int64 `query:"user_id"`
	Offset int32 `query:"offset"`
	Limit  int32 `query:"limit"`
}

func ListPublishedVideo(c *fiber.Ctx) error {
	var req listRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.InvalidParams)
	}
	loginUserID := c.Locals(constant.UserID).(int64)
	videoResp, err := rpc.VideoClient.Search(c.UserContext(), &pbvideo.SearchRequest{
		UserId:      req.UserID,
		LoginUserId: loginUserID,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.ServerInternal)
	}
	if len(videoResp.VideoList) <= 0 {
		return c.JSON(response.BuildVideoList(nil))
	}
	return c.JSON(response.BuildVideoList(videoResp.GetVideoList()))
}
