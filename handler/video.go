package handler

import (
	"douyin/package/util"
	"douyin/response"
	"douyin/service"

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
		claims, err := util.ParseToken(feedService.Token)
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
