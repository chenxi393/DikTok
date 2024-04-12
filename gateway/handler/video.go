package handler

import (
	"douyin/gateway/response"
	pbvideo "douyin/grpc/video"
	"douyin/package/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	VideoClient pbvideo.VideoClient
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
		res := response.FeedResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	res, err := VideoClient.Feed(c.UserContext(), &pbvideo.FeedRequest{
		LatestTime: req.LatestTime,
		Topic:      req.Topic,
		UserID:     userID,
	})
	if err != nil {
		res := response.FeedResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

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
	res, err := VideoClient.Search(c.UserContext(), &pbvideo.SearchRequest{
		Keyword: req.Keyword,
		UserID:  userID,
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
