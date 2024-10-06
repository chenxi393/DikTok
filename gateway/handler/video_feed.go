package handler

import (
	"context"
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
		return c.JSON(constant.InvalidParams)
	}
	userID := c.Locals(constant.UserID).(int64)
	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()
	res, err := rpc.VideoClient.Feed(ctx, &pbvideo.FeedRequest{
		LatestTime:  req.LatestTime,
		Topic:       req.Topic,
		LoginUserId: userID,
	})
	if err != nil {
		return c.JSON(constant.ServerInternal)
	}
	return c.JSON(res)
}
