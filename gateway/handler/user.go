package handler

import (
	"douyin/gateway/auth"
	"douyin/gateway/response"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"

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
		zap.L().Error(err.Error())
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
	res, err := UserClient.Login(c.UserContext(), &pbuser.LoginRequest{
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
	updateAvatar     = 3
	updateBackground = 4
	updateSignature  = 5
)

func UserUpdate(c *fiber.Ctx) error {
	// 	var req updateRequest
	// 	err := c.BodyParser(req)
	// 	if err != nil {
	// 		otelzap.Ctx(c.UserContext()).Error(err.Error())
	// 		zap.L().Error(err.Error())
	// 		res := response.UserRegisterOrLogin{
	// 			StatusCode: constant.Failed,
	// 			StatusMsg:  constant.BadParaRequest,
	// 		}
	// 		return c.JSON(res)
	// 	}
	// 	switch req.UpdateType {
	// 	case updateUsername:
	// 		avatarHeader, err := c.FormFile("avatar")
	// 		if err != nil && err != fasthttp.ErrMissingFile {
	// 			zap.L().Error(err.Error())
	// 			res := response.CommonResponse{
	// 				StatusCode: constant.Failed,
	// 				StatusMsg:  constant.FileFormatError,
	// 			}
	// 			return c.JSON(res)
	// 		}
	// 	case updatePassword:
	// 	case updateAvatar:
	// 	case updateBackground:
	// 	case updateSignature:
	// 	}

	// 	backgroundHeader, err := c.FormFile("background_image")
	// 	if err != nil && err != fasthttp.ErrMissingFile {
	// 		zap.L().Error(err.Error())
	// 		res := response.CommonResponse{
	// 			StatusCode: constant.Failed,
	// 			StatusMsg:  constant.FileFormatError,
	// 		}
	// 		return c.JSON(res)
	// 	}
	// 	res, err := UserClient.Update(c.UserContext(), &pbuser.UpdateRequest{
	// 		Username:    req.Username,
	// 		OldPassword: req.OldPassword,
	// 		NewPassword: req.NewPassword,
	// 		UserID:      c.Locals(constant.UserID).(uint64),
	// 		Signature:   req.Signature,
	// 	})
	// 	if err != nil {
	// 		res := response.UserRegisterOrLogin{
	// 			StatusCode: constant.Failed,
	// 			StatusMsg:  err.Error(),
	// 		}
	// 		return c.JSON(res)
	// 	}
	// return c.JSON(res)
	return nil
}
