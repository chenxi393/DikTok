package main

import (
	"douyin/config"
	"douyin/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func initFiber() {
	app := fiber.New()
	app.Use(logger.New()) // 使用中间件打印日志
	initRouter(app)
	panic(app.Listen(
		config.SystemConfig.HttpAddress.Host + ":" + config.SystemConfig.HttpAddress.Port).Error())
}

func initRouter(app *fiber.App) {
	app.Static("/static", "./public")
	//用户登录数据保存在内存中，单次运行过程中有效
	//视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

	api := app.Group("/douyin")
	{
		// 客户端（前端） 用户注册或者登录后 紧接着就调用 /douyin/user/
		user := api.Group("/user")
		//user.Get("/", controller.UserInfo)
		{
			user.Post("/register/", controller.UserRegister)
			user.Post("/login/", controller.UserLogin)
		}

	}

	// apiRouter.GET("/feed/", controller.Feed)
	// apiRouter.GET("/user/", controller.UserInfo)
	// apiRouter.POST("/user/register/", controller.Register)
	// apiRouter.POST("/user/login/", controller.Login)
	// apiRouter.POST("/publish/action/", controller.Publish)
	// apiRouter.GET("/publish/list/", controller.PublishList)

	// // extra apis - I
	// apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	// apiRouter.GET("/favorite/list/", controller.FavoriteList)
	// apiRouter.POST("/comment/action/", controller.CommentAction)
	// apiRouter.GET("/comment/list/", controller.CommentList)

	// // extra apis - II
	// apiRouter.POST("/relation/action/", controller.RelationAction)
	// apiRouter.GET("/relation/follow/list/", controller.FollowList)
	// apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	// apiRouter.GET("/relation/friend/list/", controller.FriendList)
	// apiRouter.GET("/message/chat/", controller.MessageChat)
	// apiRouter.POST("/message/action/", controller.MessageAction)
}
