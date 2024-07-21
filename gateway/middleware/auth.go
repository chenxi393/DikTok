package middleware

import (
	"fmt"
	"time"

	"diktok/config"
	"diktok/gateway/response"
	"diktok/package/constant"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UserClaims struct {
	UserID int64
	jwt.RegisteredClaims
}

func SignToken(userID int64) (string, error) {
	signingKey := []byte(config.System.JwtSecret)
	// 配置 userClaims ,并生成 token
	claims := UserClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constant.TokenExpiration)),
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
	token := getToken(c)
	if token == "" {
		zap.L().Info("token为空")
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.WrongToken,
		}
		return c.JSON(res)
	}
	claims, err := ParseToken(token)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.WrongToken,
		}
		return c.JSON(res)
	}
	refreshToken(c, claims)
	c.Locals(constant.UserID, claims.UserID)
	return c.Next()
}

func refreshToken(c *fiber.Ctx, claims *UserClaims) {
	// 如果过期时间 小于1天 则续期
	if claims.ExpiresAt.Before(time.Now().AddDate(0, 0, 1)) {
		t, _ := SignToken(claims.UserID)
		SetTokenCookie(c, t)
	}
}

func SetTokenCookie(c *fiber.Ctx, token string) {
	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.Expires = time.Now().Add(constant.TokenExpiration)

	// Set cookie
	c.Cookie(cookie)
}

func AuthenticationOption(c *fiber.Ctx) error {
	if token := getToken(c); token != "" {
		claims, err := ParseToken(token)
		if err != nil {
			res := response.CommonResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.WrongToken,
			}
			return c.JSON(res)
		}
		refreshToken(c, claims)
		c.Locals(constant.UserID, claims.UserID)
	} else {
		// 解决 c.locals 无数据 反射panic的问题
		c.Locals(constant.UserID, int64(0))
	}
	return c.Next()
}

func getToken(c *fiber.Ctx) (token string) {
	// 兼容 三种传参
	token = c.Get("token")
	if token == "" {
		token = c.Cookies("token")
	}
	if token == "" {
		token = c.Query("token")
	}
	return
}
