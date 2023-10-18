package router

import (
	"douyin/handler"

	"github.com/gofiber/fiber/v2"
)

func InitRouter(app *fiber.App) {
	//用户登录数据保存在内存中，单次运行过程中有效
	//视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可
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
		comment := api.Group("/comment")
		{
			comment.Post("/action/", handler.CommentAction)
			comment.Get("/list/", handler.CommentList)
		}
		relation := api.Group("/relation")
		{
			relation.Post("/action/", handler.RelationAction)
			relation.Get("/follow/list/", handler.FollowList)
			relation.Get("/follower/list/", handler.FollowerList)
			relation.Get("/friend/list/", handler.FriendList)
		}
		messgae := api.Group("/message")
		{
			messgae.Post("/action/", handler.MessageAction)
			messgae.Get("/chat/", handler.MessageChat)
		}

	}
}
