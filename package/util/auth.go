package util

import (
	"douyin/config"
	"douyin/package/constant"
	"douyin/response"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UserClaims struct {
	UserID uint64
	jwt.RegisteredClaims
}

func SignToken(userID uint64) (string, error) {
	signingKey := []byte(config.System.JwtSecret)
	// 配置 userClaims ,并生成 token
	claims := UserClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constant.TokenTimeOut)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ParseToken(token string) (*UserClaims, error) {
	signingKey := []byte(config.System.JwtSecret)
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		zap.L().Info("util.jwt.ParseToken err:", zap.Error(err))
		return nil, err
	} else if tokenClaims == nil {
		err = fmt.Errorf("jwtToken身份识别失败")
		zap.L().Info("util.jwt.ParseToken err:", zap.Error(err))
		return nil, err
	}
	if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
		return claims, nil
	}
	err = fmt.Errorf("jwtToken身份识别失败")
	zap.L().Info("util.jwt.ParseToken err:", zap.Error(err))
	return nil, err
}

func Authentication(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		zap.L().Info("token为空")
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		return c.JSON(res)
	}
	claims, err := ParseToken(token)
	if err != nil {
		res := response.UserRegisterOrLogin{
			StatusCode: response.Failed,
			StatusMsg:  response.WrongToken,
		}
		return c.JSON(res)
	}
	c.Locals(constant.UserID, claims.UserID)
	return c.Next()
}
