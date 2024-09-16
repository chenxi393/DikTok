package handler

import (
	"diktok/gateway/middleware"
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
		return c.JSON(constant.InvalidParams)
	}
	var loginUserID int64
	if req.Token == "" {
		loginUserID = 0
	} else {
		claims, err := middleware.ParseToken(req.Token)
		if err != nil {
			return c.JSON(constant.InvalidToken)
		}
		loginUserID = claims.UserID
	}
	res, err := rpc.UserClient.List(c.UserContext(), &pbuser.ListReq{
		UserID:      []int64{req.UserID},
		LoginUserID: loginUserID,
	})
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}
	return c.JSON(res)
}
