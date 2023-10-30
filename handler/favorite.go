package handler

import (
	"douyin/package/constant"
	"douyin/response"
	"douyin/service"
	"fmt"

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
	userID := c.Locals(constant.UserID).(uint64)
	var resp *response.CommonResponse
	if service.ActionType == constant.DoAction {
		resp, err = service.Favorite(userID)
	} else if service.ActionType == constant.UndoAction {
		resp, err = service.UnFavorite(userID)
	} else {
		err = fmt.Errorf(constant.BadParaRequest)
	}
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
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
	userID := c.Locals(constant.UserID).(uint64)
	resp, err := service.FavoriteList(userID)
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
