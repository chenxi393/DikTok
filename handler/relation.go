package handler

import (
	"douyin/package/util"
	"douyin/response"
	"douyin/service"

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
	claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	err = service.RelationAction(claims.UserID)
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
		StatusMsg:  "关注成功",
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
	claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := service.RelationFollowList(claims.UserID)
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
	claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := service.RelationFollowerList(claims.UserID)
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
	claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.RelationListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	if claims.UserID != service.UserID {
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
