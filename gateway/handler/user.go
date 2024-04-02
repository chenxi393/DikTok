package handler

import (
	"bytes"
	"context"
	"douyin/gateway/auth"
	"douyin/gateway/response"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"io"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var (
	UserClient pbuser.UserClient
)

type userRequest struct {
	// 密码，最长32个字符
	Password string `query:"password"`
	// 注册用户名，最长32个字符
	Username string `query:"username"`
	// 用户鉴权token
	Token string `query:"token"`
	// 用户id 注意上面token会带一个userID
	UserID uint64 `query:"user_id"`
}

func UserRegister(c *fiber.Ctx) error {
	var req userRequest
	err := c.QueryParser(&req)
	if err != nil {
		otelzap.Ctx(c.UserContext()).Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res, err := UserClient.Register(c.UserContext(), &pbuser.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	// 签发token
	token, err := auth.SignToken(res.UserId)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res.Token = token
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

func UserLogin(c *fiber.Ctx) error {
	var req userRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	// 初始化一个带取消功能的ctx 超时控制 ！ TODO 超时控制
	// 注意这里的dial 
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()
	res, err := UserClient.Login(ctx, &pbuser.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	// 签发token
	token, err := auth.SignToken(uint64(res.UserId))
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res.Token = token
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

func UserInfo(c *fiber.Ctx) error {
	var req userRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var loginUserID uint64
	if req.Token == "" {
		loginUserID = 0
	} else {
		claims, err := auth.ParseToken(req.Token)
		if err != nil {
			res := response.UserRegisterOrLogin{
				StatusCode: constant.Failed,
				StatusMsg:  constant.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		loginUserID = claims.UserID
	}
	res, err := UserClient.Info(c.UserContext(), &pbuser.InfoRequest{
		UserID:      req.UserID,
		LoginUserID: loginUserID,
	})
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

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
	userID := c.Locals(constant.UserID).(uint64)
	switch req.UpdateType {
	case updateUsername, updatePassword:
		{
			updateRes, err = UserClient.Update(c.UserContext(), &pbuser.UpdateRequest{
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
				res := response.CommonResponse{
					StatusCode: constant.Failed,
					StatusMsg:  constant.TooLongSignature,
				}
				return c.JSON(res)
			}
			updateRes, err = UserClient.Update(c.UserContext(), &pbuser.UpdateRequest{
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
				res := response.CommonResponse{
					StatusCode: constant.Failed,
					StatusMsg:  constant.FileFormatError,
				}
				return c.JSON(res)
			}
			otelzap.L().Ctx(c.UserContext()).Info("[UserUpdate] Filename:" + fileHeader.Filename)
			file, err = fileHeader.Open()
			if err != nil {
				otelzap.L().Error(err.Error())
				res := response.CommonResponse{
					StatusCode: constant.Failed,
					StatusMsg:  constant.FileFormatError,
				}
				return c.JSON(res)
			}
			defer file.Close()
			buf := bytes.NewBuffer(nil)
			if _, err = io.Copy(buf, file); err != nil {
				otelzap.L().Error(err.Error())
				res := response.CommonResponse{
					StatusCode: constant.Failed,
					StatusMsg:  constant.FileFormatError,
				}
				return c.JSON(res)
			}
			updateRes, err = UserClient.Update(c.UserContext(), &pbuser.UpdateRequest{
				UserID:     userID,
				UpdateType: req.UpdateType,
				Data:       buf.Bytes(),
			})
		}
	}
	if err != nil {
		otelzap.L().Ctx(c.UserContext()).Error(err.Error())
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		return c.JSON(res)
	}
	return c.JSON(updateRes)
}

type userSearchReq struct {
	username string `query:"username"`
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
	// userID := c.Locals(constant.UserID).(uint64)
	// var updateRes *pbuser.UpdateResponse
	// var fileHeader *multipart.FileHeader
	// var file multipart.File
}
