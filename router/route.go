package router

import (
	"douyin/handler"
	"douyin/package/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouter(app *fiber.App) {
	app.Use(cors.New())
	app.Static("/video", "./douyinVideo",
		fiber.Static{ByteRange: true}) // 好像可以分块传输 但是客户端没啥用。
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
			publish.Get("/list/", util.Authentication, handler.ListPublishedVideo)
		}
		favorite := api.Group("/favorite", util.Authentication)
		{
			favorite.Post("/action/", handler.FavoriteVideoAction)
			favorite.Get("/list/", handler.FavoriteList)
		}
		comment := api.Group("/comment")
		{
			comment.Post("/action/", util.Authentication, handler.CommentAction)
			comment.Get("/list/", handler.CommentList)
		}
		relation := api.Group("/relation", util.Authentication) // 这里暂时关注和粉丝列表都需要鉴权
		{
			relation.Post("/action/", handler.RelationAction)
			relation.Get("/follow/list/", handler.FollowList)
			relation.Get("/follower/list/", handler.FollowerList)
			relation.Get("/friend/list/", handler.FriendList)
		}
		messgae := api.Group("/message")
		{
			messgae.Post("/action/", util.Authentication, handler.MessageAction)
			messgae.Get("/chat/", util.Authentication, handler.MessageChat)
		}
	}
}
