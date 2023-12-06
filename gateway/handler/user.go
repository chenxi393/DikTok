package handler

import (
	"context"
	"douyin/gateway/auth"
	pbuser "douyin/grpc/user"
	"douyin/response"

	"github.com/gofiber/fiber/v2"
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
	var user userRequest
	err := c.QueryParser(&user)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res, err := UserClient.Register(context.Background(), &pbuser.RegisterRequest{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
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
			StatusCode: response.Failed,
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
	var user userRequest
	err := c.QueryParser(&user)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res, err := UserClient.Login(context.Background(), &pbuser.LoginRequest{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
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
			StatusCode: response.Failed,
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
	var user userRequest
	err := c.QueryParser(&user)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var loginUserID uint64
	if user.Token == "" {
		loginUserID = 0
	} else {
		claims, err := auth.ParseToken(user.Token)
		if err != nil {
			res := response.UserRegisterOrLogin{
				StatusCode: response.Failed,
				StatusMsg:  response.WrongToken,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		loginUserID = claims.UserID
	}
	res, err := UserClient.Info(context.Background(), &pbuser.InfoRequest{
		UserID:      user.UserID,
		LoginUserID: loginUserID,
	})
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
