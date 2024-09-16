package handler

import (
	"diktok/gateway/response"
	"diktok/package/constant"
	"diktok/package/util"

	"github.com/gofiber/fiber/v2"
)

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
		return c.JSON(constant.InvalidParams)
	}
	key, err := util.GetUUid()
	if err != nil {
		return c.JSON(constant.ServerInternal)
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
