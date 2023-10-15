package handler

import (
	"douyin/package/util"
	"douyin/response"
	"douyin/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func FavoriteVideoAction(c *fiber.Ctx) error {
	var service service.FavoriteService
	err := c.QueryParser(&service)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	// 鉴权
	Claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	err = service.FavoriteAction(Claims.UserID)
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
		StatusMsg:  "点赞成功",
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

func FavoriteList(c *fiber.Ctx) error {
	var service service.FavoriteService
	err := c.QueryParser(&service)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.PublishListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
			VideoList:  nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	// 鉴权
	Claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.PublishListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
			VideoList:  nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := service.FavoriteList(Claims.UserID)
	if err != nil {
		res := response.PublishListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
			VideoList:  nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res := response.PublishListResponse{
		StatusCode: response.Success,
		StatusMsg:  "查询喜欢视频列表成功",
		VideoList:  resp,
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
