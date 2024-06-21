package handler

import (
	"bytes"
	"io"
	"strings"

	"diktok/gateway/auth"
	"diktok/gateway/response"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/util"

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

type listRequest struct {
	// 用户鉴权token
	Token  string `query:"token"`
	UserID int64  `query:"user_id"`
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
	userClaim, err := auth.ParseToken(req.Token)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.WrongToken,
		}
		return c.JSON(res)
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
	res, err := VideoClient.Publish(c.UserContext(), &pbvideo.PublishRequest{
		Title:  req.Title,
		Topic:  req.Topic,
		UserID: userClaim.UserID,
		Data:   buf.Bytes(),
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

func ListPublishedVideo(c *fiber.Ctx) error {
	var req listRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.VideoListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	resp, err := VideoClient.List(c.UserContext(), &pbvideo.ListRequest{
		UserID:      req.UserID,
		LoginUserID: userID,
	})
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	return c.JSON(resp)
}

const (
	video      = "1"
	avatar     = "2"
	background = "3"
)

func UploadToken(c *fiber.Ctx) error {
	uploadType := c.Get("upload_type")
	var prefix string
	switch uploadType {
	case video:
		prefix = "V-"
	case avatar:
		prefix = "A-"
	case background:
		prefix = "B-"
	default:
		// TODO 这些错误返回可以封装一下 给一个 code msg 和 结构体 直接封装 之前的也是 学一下
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	key, err := util.GetUUid()
	if err != nil {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	key = prefix + key
	token := util.GetUploadToken(key)
	res := response.UploadTokenResponse{
		StatusCode:  constant.Success,
		StatusMsg:   constant.GetTokenSuccess,
		UploadToken: token,
		FileName:    key,
	}
	return c.JSON(res)
}
