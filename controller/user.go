package controller

import (
	"douyin/service"

	"github.com/gofiber/fiber/v2"
)

type RegisterResponse struct {
	// 状态码，0-成功，其他值-失败
	StatusCode int64 `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
	// 用户鉴权token
	Token string `json:"token"`
	// 用户id
	UserID int64 `json:"user_id"`
}

func UserRegister(c *fiber.Ctx) error {
	var userService service.UserService
	err := c.QueryParser(&userService)
	if err != nil {
		res := RegisterResponse{
			StatusCode: -1,
			StatusMsg:  "参数错误，注册失败",
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
	// 参数匹配正确开始注册
	res, err := userService.RegisterService()
	if err != nil {
		res := RegisterResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}

func UserLogin(c *fiber.Ctx) error {
	var userService service.UserService
	err := c.QueryParser(&userService)
	if err != nil {
		res := RegisterResponse{
			StatusCode: -1,
			StatusMsg:  "参数错误，登录失败",
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}

	res, err := userService.LoginService()
	if err != nil || res == nil {
		res := RegisterResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(res)
}
