package handler

import (
	"diktok/gateway/response"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type searchRequest struct {
	Keyword string `query:"keyword"`
}

// 23.11.03 新增视频搜索功能
func SearchVideo(c *fiber.Ctx) error {
	var req searchRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.VideoListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	res, err := rpc.VideoClient.Search(c.UserContext(), &pbvideo.SearchRequest{
		Keyword:     req.Keyword,
		LoginUserId: userID,
	})
	if err != nil {
		res := response.VideoListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
