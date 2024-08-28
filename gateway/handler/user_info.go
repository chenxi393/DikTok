package handler

import (
	"diktok/gateway/middleware"
	"diktok/gateway/response"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

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
	var loginUserID int64
	if req.Token == "" {
		loginUserID = 0
	} else {
		claims, err := middleware.ParseToken(req.Token)
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
	res, err := rpc.UserClient.List(c.UserContext(), &pbuser.ListReq{
		UserID:      []int64{req.UserID},
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
