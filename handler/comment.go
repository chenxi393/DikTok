package handler

import (
	"douyin/package/util"
	"douyin/response"
	"douyin/service"

	"github.com/gofiber/fiber/v2"
)

func CommentAction(c *fiber.Ctx) error {
	var service service.CommentService
	err := c.QueryParser(&service)
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		c.JSON(res)
	}
	Claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := service.CommentAction(Claims.UserID)
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func CommentList(c *fiber.Ctx) error {
	var service service.CommentService
	err := c.QueryParser(&service)
	if err != nil {
		res := response.CommentListResponse{
			StatusCode:  response.Failed,
			StatusMsg:   response.BadParaRequest,
			CommentList: nil,
		}
		c.Status(fiber.StatusOK)
		c.JSON(res)
	}
	claims, err := util.ParseToken(service.Token)
	if err != nil {
		res := response.CommentListResponse{
			StatusCode:  response.Failed,
			StatusMsg:   response.WrongToken,
			CommentList: nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	resp, err := service.CommentList(claims.UserID)
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
