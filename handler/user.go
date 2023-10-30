package handler

import (
	"douyin/package/util"
	"douyin/response"
	"douyin/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func UserRegister(c *fiber.Ctx) error {
	var userService service.UserService
	err := c.QueryParser(&userService)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res, err := userService.RegisterService()
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

func UserLogin(c *fiber.Ctx) error {
	var userService service.UserService
	err := c.QueryParser(&userService)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res, err := userService.LoginService()
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

func UserInfo(c *fiber.Ctx) error {
	var userService service.UserService
	err := c.QueryParser(&userService)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	claims, err := util.ParseToken(userService.Token)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	res, err := userService.InfoService(claims.UserID)
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
