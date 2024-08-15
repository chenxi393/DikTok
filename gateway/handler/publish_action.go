package handler

import (
	"bytes"
	"io"
	"strings"

	"diktok/gateway/middleware"
	"diktok/gateway/response"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
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
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	if userID == 0 {
		userClaim, err := middleware.ParseToken(req.Token)
		if err != nil {
			zap.L().Error(err.Error())
			res := response.CommonResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.WrongToken,
			}
			return c.JSON(res)
		}
		userID = userClaim.UserID
	}

	fileHeader, err := c.FormFile("data")
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	// 检查文件后缀是不是mp4 大小在上传的时候会限制30MB
	if !strings.HasSuffix(fileHeader.Filename, constant.MP4Suffix) {
		zap.L().Error(constant.FileFormatError)
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	zap.L().Info("PublishAction Filename:" + fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	defer file.Close()
	// 将文件转化为字节流
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.FileFormatError,
		}
		return c.JSON(res)
	}
	res, err := rpc.VideoClient.Publish(c.UserContext(), &pbvideo.PublishRequest{
		Title:       req.Title,
		Topic:       req.Topic,
		LoginUserId: userID,
		Data:        buf.Bytes(),
	})
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	return c.JSON(res)
}
