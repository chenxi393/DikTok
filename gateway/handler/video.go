package handler

import (
	"douyin/response"
	"douyin/service"
	"douyin/gateway/auth"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Feed(c *fiber.Ctx) error {
	var feedService service.FeedService
	err := c.QueryParser(&feedService)
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
	if feedService.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(feedService.Token)
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
	res, err := feedService.GetFeed(userID)
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
	var service service.SearchService
	err := c.QueryParser(&service)
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
	if service.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(service.Token)
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
	res, err := service.SearchVideo(userID)
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
