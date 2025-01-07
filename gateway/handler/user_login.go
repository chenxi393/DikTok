package handler

import (
	"context"
	"time"

	"diktok/gateway/middleware"
	"diktok/gateway/response"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

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
	// 初始化一个带取消功能的ctx 超时控制
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()
	res, err := rpc.UserClient.Login(ctx, &pbuser.LoginRequest{
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
	if res.StatusCode != 0 {
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	// 签发token
	token, err := middleware.SignToken(int64(res.UserId))
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	middleware.SetTokenCookie(c, token)
	res.Token = token
	c.Status(fiber.StatusOK)
	return c.JSON(response.BuildLoginRes(res))
}
