package handler

import (
	"bytes"
	"io"

	"diktok/gateway/middleware"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"github.com/h2non/filetype"
	"go.uber.org/zap"
)

type publishRequest struct {
	// 用户鉴权token
	Token string `form:"token"`
	// 视频标题
	Title string `form:"title"`
	// 新增 topic
	Topic string `form:"topic"`
}

func PublishAction(c *fiber.Ctx) error {
	var req publishRequest
	err := c.BodyParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.InvalidParams)
	}
	userID := c.Locals(constant.UserID).(int64)
	if userID == 0 {
		userClaim, err := middleware.ParseToken(req.Token)
		if err != nil {
			zap.L().Error(err.Error())
			return c.JSON(constant.InvalidToken)
		}
		userID = userClaim.UserID
	}

	fileHeader, err := c.FormFile("data")
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.FileUploadFailed)
	}

	zap.L().Info("PublishAction Filename:" + fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.FileUploadFailed)
	}
	defer file.Close()
	// 将文件转化为字节流
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.FileUploadFailed)
	}

	// 检查文件是不是mp4 大小在上传的时候会限制30MB
	if !filetype.IsVideo(buf.Bytes()) {
		return c.JSON(constant.FileIsNotVideo)
	}

	res, err := rpc.VideoClient.Publish(c.UserContext(), &pbvideo.PublishRequest{
		Title:       req.Title,
		Topic:       req.Topic,
		LoginUserId: userID,
		Data:        buf.Bytes(),
	})
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.ServerInternal)
	}
	return c.JSON(res)
}
