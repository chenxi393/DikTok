package main

import (
	"diktok/config"
	"diktok/gateway/handler"
	"diktok/gateway/middleware"
	"diktok/gateway/response"

	// "github.com/gofiber/contrib/otelfiber/v2"
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
		JSONEncoder: response.GrpcMarshal,
	})
	// 使用中间件打印日志
	app.Use(logger.New())
	initRouter(app)
	// if config.System.Mode != constant.DebugMode {
	// 	zap.L().Fatal("fiber启动失败: ", zap.Error(app.ListenTLS(
	// 		config.System.HTTP.Host+":"+config.System.HTTP.Port, "./server.crt", "./server.key")))
	// }
	zap.L().Fatal("fiber启动失败: ", zap.Error(app.Listen(
		config.System.HTTP.Host+":"+config.System.HTTP.Port)))
}

func initRouter(app *fiber.App) {
	// https://github.com/gofiber/contrib/issues/1126
	// FIXME 这都是啥bug啊 使用这个导致sse 阻塞不可用
	// app.Use(otelfiber.Middleware())

	app.Use(middleware.Recovery)
	// 允许跨域请求
	app.Use(cors.New())
	api := app.Group("/douyin")
	{
		api.Get("/feed", middleware.AuthenticationOption, handler.Feed)
		// get token of file upload
		api.Get("/upload/get_token", middleware.Authentication, handler.UploadToken)

		// 视频搜索 可以拓展搜索用户
		search := api.Group("/search", middleware.AuthenticationOption)
		{
			search.Get("/video", handler.SearchVideo)
			search.Get("/user", handler.UserSearch)
		}

		user := api.Group("/user")
		{
			user.Get("/", middleware.AuthenticationOption, handler.UserInfo)
			user.Post("/register", handler.UserRegister)
			user.Post("/login", handler.UserLogin)
			user.Post("/update", middleware.Authentication, handler.UserUpdate)
		}

		publish := api.Group("/publish")
		{
			// action token放在body端 先使用中间件鉴权 userid 无在拿body
			publish.Post("/action", middleware.AuthenticationOption, handler.PublishAction)
			publish.Get("/list", middleware.AuthenticationOption, handler.ListPublishedVideo)
		}

		favorite := api.Group("/favorite")
		{
			favorite.Post("/action", middleware.Authentication, handler.FavoriteVideoAction)
			favorite.Get("/list", middleware.AuthenticationOption, handler.FavoriteList)
		}

		comment := api.Group("/comment")
		{
			comment.Post("/action", middleware.Authentication, handler.CommentAction)
			comment.Get("/list", middleware.AuthenticationOption, handler.CommentList)
		}

		relation := api.Group("/relation")
		{
			relation.Post("/action", middleware.Authentication, handler.RelationAction)
			relation.Get("/follow/list", middleware.AuthenticationOption, handler.FollowList)
			relation.Get("/follower/list", middleware.AuthenticationOption, handler.FollowerList)
			relation.Get("/friend/list", middleware.Authentication, handler.FriendList)
		}

		messgae := api.Group("/message", middleware.Authentication)
		{
			messgae.Post("/action", handler.MessageAction)
			messgae.Get("/chat", handler.MessageChat)
			// 使用websocket替换http每秒轮询
			messgae.Use("/ws", func(c *fiber.Ctx) error {
				if websocket.IsWebSocketUpgrade(c) {
					return c.Next()
				}
				return fiber.ErrUpgradeRequired
			})
			messgae.Get("/ws", websocket.New(handler.MessageWebsocket()))
			// sse
			messgae.Post("/sse", handler.SSEHandle)
		}
	}
}
