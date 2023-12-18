package main

import (
	"douyin/config"
	"douyin/gateway/auth"
	"douyin/gateway/handler"
	"douyin/gateway/util"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

func startFiber() {
	// 客户端文件超过30MB 返回413
	app := fiber.New(fiber.Config{
		BodyLimit:   30 * 1024 * 1024,
		JSONEncoder: util.GrpcMarshal,
	})
	// 使用中间件打印日志
	app.Use(logger.New())
	initRouter(app)
	zap.L().Fatal("fiber启动失败: ", zap.Error(app.Listen(
		config.System.HttpAddress.Host+":"+config.System.HttpAddress.Port)))
}

func initRouter(app *fiber.App) {
	// 允许跨域请求
	app.Use(cors.New())
	api := app.Group("/douyin")
	{
		// 新增接口 搜索功能 可以拓展搜索用户
		search := api.Group("/search")
		{
			search.Get("/video/", handler.SearchVideo)
		}
		api.Get("/feed/", handler.Feed)
		user := api.Group("/user")
		user.Get("/", handler.UserInfo)
		{
			user.Post("/register/", handler.UserRegister)
			user.Post("/login/", handler.UserLogin)
		}
		publish := api.Group("/publish")
		{
			// action token放在body端 不适用中间件鉴权
			publish.Post("/action/", handler.PublishAction)
			publish.Get("/list/", handler.ListPublishedVideo)
		}
		favorite := api.Group("/favorite")
		{
			favorite.Post("/action/", auth.Authentication, handler.FavoriteVideoAction)
			favorite.Get("/list/", handler.FavoriteList)
		}
		comment := api.Group("/comment")
		{
			comment.Post("/action/", auth.Authentication, handler.CommentAction)
			comment.Get("/list/", handler.CommentList)
		}
		relation := api.Group("/relation")
		{
			relation.Post("/action/", auth.Authentication, handler.RelationAction)
			relation.Get("/follow/list/", handler.FollowList)
			relation.Get("/follower/list/", handler.FollowerList)
			relation.Get("/friend/list/", auth.Authentication, handler.FriendList)
		}
		messgae := api.Group("/message", auth.Authentication)
		{
			messgae.Post("/action/", handler.MessageAction)
			messgae.Get("/chat/", handler.MessageChat)
		}
		// 使用websocket替换http每秒轮询
		messgae.Use("/ws", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})
		messgae.Get("/ws", websocket.New(handleWebSocket()))
	}
}
