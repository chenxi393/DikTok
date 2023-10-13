package handler

import (
	"douyin/response"
	"douyin/service"

	"github.com/gofiber/fiber/v2"
)

func Feed(c *fiber.Ctx) error {
	var feedService service.FeedService
	err := c.QueryParser(&feedService)
	if err != nil {
		res := response.FeedResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res, err := feedService.GetFeed()
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
