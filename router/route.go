package router

import (
	"douyin/handler"

	"github.com/gofiber/fiber/v2"
)

func InitRouter(app *fiber.App) {
	//用户登录数据保存在内存中，单次运行过程中有效
	//视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可
	// jwtFunc := jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{Key: []byte(config.SystemConfig.JwtSecret)},
	// 	Claims: jwt.RegisteredClaims{},
	// })
	app.Static("/video", "./douyinVideo") // 是可以用绝对路径
	app.Static("/image", "./douyinImage") // 是可以用绝对路径
	api := app.Group("/douyin")
	{
		api.Get("/feed/", handler.Feed)

		// 客户端（前端） 用户注册或者登录后 紧接着就调用 /douyin/user/
		user := api.Group("/user")
		user.Get("/", handler.UserInfo)
		{
			user.Post("/register/", handler.UserRegister)
			user.Post("/login/", handler.UserLogin)
		}

		publish := api.Group("/publish")
		{
			publish.Post("/action/", handler.PublishAction)
			publish.Get("/list/", handler.ListPublishedVideo)
		}
		favorite := api.Group("/favorite")
		{
			favorite.Post("/action/", handler.FavoriteVideoAction)
			favorite.Get("/list/", handler.FavoriteList)
		}

	}

	// apiRouter.POST("/publish/action/", handler.Publish)
	// apiRouter.GET("/publish/list/", handler.PublishList)

	// // extra apis - I
	// apiRouter.POST("/favorite/action/", handler.FavoriteAction)
	// apiRouter.GET("/favorite/list/", handler.FavoriteList)
	// apiRouter.POST("/comment/action/", handler.CommentAction)
	// apiRouter.GET("/comment/list/", handler.CommentList)

	// // extra apis - II
	// apiRouter.POST("/relation/action/", handler.RelationAction)
	// apiRouter.GET("/relation/follow/list/", handler.FollowList)
	// apiRouter.GET("/relation/follower/list/", handler.FollowerList)
	// apiRouter.GET("/relation/friend/list/", handler.FriendList)
	// apiRouter.GET("/message/chat/", handler.MessageChat)
	// apiRouter.POST("/message/action/", handler.MessageAction)
}
