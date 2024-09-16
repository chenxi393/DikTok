package middleware

import (
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Recovery(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Sugar().Errorf("catch error: %#v", r)
			c.JSON(constant.ServerInternal)
		}
	}()
	return c.Next()
}
