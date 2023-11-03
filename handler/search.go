package handler

import (
	"douyin/package/util"
	"douyin/response"
	"douyin/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// 23.11.03 新增视频搜索功能
func SearchVideo(c *fiber.Ctx) error {
	var service service.SearchService
	err := c.QueryParser(&service)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.VideoListResponse{
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
			res := response.VideoListResponse{
				StatusCode: response.Failed,
				StatusMsg:  response.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		userID = claims.UserID
	}
	res, err := service.SearchVideo(userID)
	if err != nil {
		res := response.VideoListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
