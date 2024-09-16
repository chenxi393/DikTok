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
	videoResp, err := rpc.VideoClient.MGet(c.UserContext(), &pbvideo.MGetReq{
		UserId: req.UserID,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.ServerInternal)
	}
	data, err := BuildVideosInfo(c.Context(), nil, videoResp.VideoList, loginUserID)
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.ServerInternal)
	}
	return c.JSON(&response.VideoListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.PublishListSuccess,
		VideoList:  data,
	})
}
