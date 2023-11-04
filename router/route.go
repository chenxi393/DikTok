package router

import (
	"douyin/handler"
	"douyin/package/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouter(app *fiber.App) {
	// 允许所有跨域请求
	app.Use(cors.New())
	app.Static("/video", "./douyinVideo",
		fiber.Static{ByteRange: true}) // 好像可以分块传输 但是客户端没啥用。
	app.Static("/image", "./douyinImage") // 是可以用绝对路径

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
			favorite.Post("/action/", util.Authentication, handler.FavoriteVideoAction)
			favorite.Get("/list/", handler.FavoriteList)
		}
		comment := api.Group("/comment")
		{
			comment.Post("/action/", util.Authentication, handler.CommentAction)
			comment.Get("/list/", handler.CommentList)
		}
		relation := api.Group("/relation")
		{
			relation.Post("/action/", util.Authentication, handler.RelationAction)
			relation.Get("/follow/list/", handler.FollowList)
			relation.Get("/follower/list/", handler.FollowerList)
			relation.Get("/friend/list/", util.Authentication, handler.FriendList)
		}
		messgae := api.Group("/message")
		{
			messgae.Post("/action/", util.Authentication, handler.MessageAction)
			messgae.Get("/chat/", util.Authentication, handler.MessageChat)
		}
	}
}
