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
		return c.JSON(response.BuildStdResp(constant.Failed, constant.BadParaRequest, nil))
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
