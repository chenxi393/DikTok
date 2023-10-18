package handler

import (
	"douyin/package/util"
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
	claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	err = service.MessageAction(claims.UserID)
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
		StatusMsg:  "发送成功",
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
	claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.MessageResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := service.MessageChat(claims.UserID)
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
