package handler

import (
	"douyin/package/constant"
	"douyin/package/util"
	"douyin/response"
	"douyin/service"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RelationAction(c *fiber.Ctx) error {
	var service service.RelationService
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
	if service.ActionType == "1" {
		err = service.FollowAction(userID)
	} else if service.ActionType == "2" {
		err = service.UnFollowAction(userID)
	} else {
		err = errors.New(constant.BadParaRequest)
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
	res := response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  response.ActionSuccess,
	}
	return c.JSON(res)
}

func FollowList(c *fiber.Ctx) error {
	var service service.RelationService
	err := c.QueryParser(&service)
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
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
	resp, err := service.RelationFollowList(userID)
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func FollowerList(c *fiber.Ctx) error {
	var service service.RelationService
	err := c.QueryParser(&service)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.RelationListResponse{
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
	resp, err := service.RelationFollowerList(userID)
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func FriendList(c *fiber.Ctx) error {
	var service service.RelationService
	err := c.QueryParser(&service)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	if userID != service.UserID {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.FriendListError,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := service.RelationFriendList()
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
