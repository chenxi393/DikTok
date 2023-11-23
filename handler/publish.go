package handler

import (
	"bytes"
	"douyin/package/constant"
	"douyin/package/util"
	"douyin/response"
	"douyin/service"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func PublishAction(c *fiber.Ctx) error {
	var publishService service.PublisService
	err := c.BodyParser(&publishService)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		return c.JSON(res)
	}
	userClaim, err := util.ParseToken(publishService.Token)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		return c.JSON(res)
	}
	fileHeader, err := c.FormFile("data")
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.FileFormatError,
		}
		return c.JSON(res)
	}
	// 检查文件后缀是不是mp4 大小在上传的时候会限制30MB
	if !strings.HasSuffix(fileHeader.Filename, constant.MP4Suffix) {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.FileFormatError,
		}
		return c.JSON(res)
	}
	zap.L().Info("PublishAction Filename:" + fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.FileFormatError,
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
			StatusMsg:  response.FileFormatError,
		}
		return c.JSON(res)
	}
	res, err := publishService.PublishAction(userClaim.UserID, buf)
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
		res := response.VideoListResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		return c.JSON(res)
	}
	var userID uint64
	if listService.Token == "" {
		userID = 0
	} else {
		claims, err := util.ParseToken(listService.Token)
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
	resp, err := listService.GetPublishVideos(userID)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	return c.JSON(resp)
}
