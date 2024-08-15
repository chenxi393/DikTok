package handler

import (
	"context"
	"diktok/gateway/response"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type feedRequest struct {
	LatestTime int64 `query:"latest_time"`
	// 新增topic
	Topic string `query:"topic"`
}

func Feed(c *fiber.Ctx) error {
	var req feedRequest
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
	ctx, cancel := context.WithTimeout(c.UserContext(), time.Second)
	defer cancel()
	res, err := rpc.VideoClient.Feed(ctx, &pbvideo.FeedRequest{
		LatestTime:  req.LatestTime,
		Topic:       req.Topic,
		LoginUserId: userID,
	})
	if err != nil {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
