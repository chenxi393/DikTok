package handler

import (
	"diktok/gateway/response"
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type userSearchReq struct {
	Username string `query:"username"`
}

func UserSearch(c *fiber.Ctx) error {
	var req userSearchReq
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	return nil
	// userID := c.Locals(constant.UserID).(int64)
	// var updateRes *pbuser.UpdateResponse
	// var fileHeader *multipart.FileHeader
	// var file multipart.File
}
