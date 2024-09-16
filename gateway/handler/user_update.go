package handler

import (
	"bytes"
	"io"
	"mime/multipart"

	"diktok/gateway/response"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type updateRequest struct {
	// 注册用户名，最长32个字符
	Username string `form:"username"`
	// 密码，最长32个字符
	OldPassword string `form:"old_password"`
	NewPassword string `form:"new_password"`
	Signature   string `form:"signature"`
	UpdateType  int32  `form:"update_type"`
}

const (
	updateUsername   = 1
	updatePassword   = 2
	updateSignature  = 3
	updateAvatar     = 4
	updateBackground = 5
)

func UserUpdate(c *fiber.Ctx) error {
	var req updateRequest
	err := c.BodyParser(&req)
	if err != nil {
		otelzap.L().Ctx(c.UserContext()).Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		return c.JSON(res)
	}
	var updateRes *pbuser.UpdateResponse
	var fileHeader *multipart.FileHeader
	var file multipart.File
	userID := c.Locals(constant.UserID).(int64)
	switch req.UpdateType {
	case updateUsername, updatePassword:
		{
			updateRes, err = rpc.UserClient.Update(c.UserContext(), &pbuser.UpdateRequest{
				UpdateType:  req.UpdateType,
				Username:    req.Username,
				UserID:      userID,
				OldPassword: req.OldPassword,
				NewPassword: req.NewPassword,
			})
		}
	case updateSignature:
		{
			if req.Signature == "" || len(req.Signature) > 255 {
				return c.JSON(constant.InvalidParams.WithDetails(constant.TooLongSignature))
			}
			updateRes, err = rpc.UserClient.Update(c.UserContext(), &pbuser.UpdateRequest{
				UpdateType: req.UpdateType,
				UserID:     userID,
				Signature:  req.Signature,
			})
		}
	case updateAvatar, updateBackground:
		{
			fileHeader, err = c.FormFile("data")
			if err != nil {
				otelzap.L().Error(err.Error())
				return c.JSON(constant.FileUploadFailed)
			}
			otelzap.L().Ctx(c.UserContext()).Info("[UserUpdate] Filename:" + fileHeader.Filename)
			file, err = fileHeader.Open()
			if err != nil {
				otelzap.L().Error(err.Error())
				return c.JSON(constant.FileUploadFailed)
			}
			defer file.Close()
			buf := bytes.NewBuffer(nil)
			if _, err = io.Copy(buf, file); err != nil {
				otelzap.L().Error(err.Error())
				return c.JSON(constant.FileUploadFailed)
			}
			updateRes, err = rpc.UserClient.Update(c.UserContext(), &pbuser.UpdateRequest{
				UserID:     userID,
				UpdateType: req.UpdateType,
				Data:       buf.Bytes(),
			})
		}
	}
	if err != nil {
		otelzap.L().Ctx(c.UserContext()).Error(err.Error())
		return c.JSON(constant.ServerInternal)
	}
	return c.JSON(updateRes)
}
