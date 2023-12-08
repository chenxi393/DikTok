package handler

import (
	"context"
	"douyin/gateway/auth"
	pbvideo "douyin/grpc/video"
	"douyin/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	VideoClient pbvideo.VideoClient
)

type feedRequest struct {
	LatestTime int64 `query:"latest_time"`
	// 用户登录状态下设置
	Token string `query:"token"`
	// 新增topic
	Topic string `query:"topic"`
}

type searchRequest struct {
	Keyword string `query:"keyword"`
	// 用户登录状态下设置
	Token string `query:"token"`
}

func Feed(c *fiber.Ctx) error {
	var req feedRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.FeedResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var userID uint64
	if req.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(req.Token)
		if err != nil {
			res := response.UserRegisterOrLogin{
				StatusCode: response.Failed,
				StatusMsg:  response.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		userID = claims.UserID
	}
	res, err := VideoClient.Feed(context.Background(), &pbvideo.FeedRequest{
		LatestTime: req.LatestTime,
		Topic:      req.Topic,
		UserID:     userID,
	})
	if err != nil {
		res := response.FeedResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

// 23.11.03 新增视频搜索功能
func SearchVideo(c *fiber.Ctx) error {
	var req searchRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.VideoListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var userID uint64
	if req.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(req.Token)
		if err != nil {
			res := response.VideoListResponse{
				StatusCode: response.Failed,
				StatusMsg:  response.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		userID = claims.UserID
	}
	res, err := VideoClient.Search(context.Background(), &pbvideo.SearchRequest{
		Keyword: req.Keyword,
		UserID:  userID,
	})
	if err != nil {
		res := response.VideoListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
