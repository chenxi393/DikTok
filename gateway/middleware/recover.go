package middleware

import (
	"diktok/gateway/response"
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recovery(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Sugar().Errorf("catch error: %#v", r)
			res := response.CommonResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.InternalException,
			}
			c.JSON(res)
		}
	}()
	return c.Next()
}
