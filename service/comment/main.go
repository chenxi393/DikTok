package main

import (
	"diktok/config"
	pbcomment "diktok/grpc/comment"
	"diktok/package/constant"
	"diktok/package/nacos"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/comment/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	nacos.InitNacos()
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://comment", constant.ServiceName+".comment")
	// defer shutdown()
	database.InitMySQL()
	storage.CommentRedis = cache.InitRedis(config.System.Redis.CommentDB)
	ConnClose := rpc.InitRpcClientWithNacos(nacos.GetNamingClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServerWithNacos(constant.CommentAddr, constant.CommentService, pbcomment.RegisterCommentServer, &CommentService{})
}
