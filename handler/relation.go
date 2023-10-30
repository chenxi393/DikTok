package handler

import (
	"douyin/package/constant"
	"douyin/response"
	"douyin/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
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
		err = fmt.Errorf("参数类型错误")
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
		StatusMsg:  "操作成功",
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
	userID := c.Locals(constant.UserID).(uint64)
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
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
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
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	if userID != service.UserID {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  "无法查看别人的好友列表",
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
