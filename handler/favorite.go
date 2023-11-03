package handler

import (
	"douyin/package/constant"
	"douyin/package/util"
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
		res := response.VideoListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
			VideoList:  nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var userID uint64
	if service.Token == "" {
		userID = 0
	} else {
		claims, err := util.ParseToken(service.Token)
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
	resp, err := service.FavoriteList(userID)
	if err != nil {
		res := response.VideoListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
			VideoList:  nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res := response.VideoListResponse{
		StatusCode: response.Success,
		StatusMsg:  response.FavoriteListSuccess,
		VideoList:  resp,
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
