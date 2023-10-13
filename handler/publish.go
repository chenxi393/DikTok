package handler

import (
	"bytes"
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
	// FIX 似乎鉴权甚至应该在参数校验之前 我记得商城是使用鉴权中间件的
	// token中间件使用失败了 目前还是手动调用  失败原因 应该还是没有找到合适的API
	Claims, err := util.ParseToken(publishService.Token)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	// TODO：怎么拿到视频数据 这一块HTTP视频传输还有一些API还是不清楚
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
	res, err := publishService.PublishAction(Claims.UserID, buf.Bytes())
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
	// 需要鉴权
	Claims, err := util.ParseToken(listService.Token)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		return c.JSON(res)
	}
	// 给出user_id 的所有视频信息
	// 并查询 Claim.user_id 看看它有没有点赞 还有有没有关注此人
	resp, err := listService.GetPublishVideos(Claims.UserID)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		return c.JSON(res)
	}
	return c.JSON(resp)
}
