package handler

import (
	"douyin/package/constant"
	"douyin/response"
	"douyin/service"

	"github.com/gofiber/fiber/v2"
)

func MessageAction(c *fiber.Ctx) error {
	var service service.MessageService
	err := c.QueryParser(&service)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	err = service.MessageAction(userID)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res := response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  response.SendSuccess,
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

// 一旦进入消息界面 客户端每秒会调用一次
func MessageChat(c *fiber.Ctx) error {
	var service service.MessageService
	err := c.QueryParser(&service)
	if err != nil {
		res := response.MessageResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	resp, err := service.MessageChat(userID)
	if err != nil {
		res := response.MessageResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
