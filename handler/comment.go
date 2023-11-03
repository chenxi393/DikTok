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

func CommentAction(c *fiber.Ctx) error {
	var service service.CommentService
	err := c.QueryParser(&service)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	var resp *response.CommentActionResponse
	if service.ActionType == constant.DoAction && service.CommentText != nil {
		resp, err = service.PostComment(userID)
	} else if service.ActionType == constant.UndoAction && service.CommentID != nil {
		resp, err = service.DeleteComment(userID)
	} else {
		err = errors.New(constant.BadParaRequest)
	}
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
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
		return c.JSON(res)
	}
	var userID uint64
	if service.Token == "" {
		userID = 0
	} else {
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
		userID = claims.UserID
	}
	resp, err := service.CommentList(userID)
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
