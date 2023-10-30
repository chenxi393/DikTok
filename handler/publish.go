package handler

import (
	"bytes"
	"douyin/package/constant"
	"douyin/package/util"
	"douyin/response"
	"douyin/service"
	"io"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func PublishAction(c *fiber.Ctx) error {
	var publishService service.PublisService
	err := c.BodyParser(&publishService)
	if err != nil {
		zap.L().Info(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		return c.JSON(res)
	}
	claims, err := util.ParseToken(publishService.Token)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	fileHeader, err := c.FormFile("data")
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	zap.L().Info("handler.publish_service.PublishAction Filename:" + fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	defer file.Close()
	// 将文件转化为字节流
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	res, err := publishService.PublishAction(claims.UserID, buf)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	return c.JSON(res)
}

func ListPublishedVideo(c *fiber.Ctx) error {
	var listService service.PublishListService
	err := c.QueryParser(&listService)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.PublishListResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	resp, err := listService.GetPublishVideos(userID)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		return c.JSON(res)
	}
	return c.JSON(resp)
}
