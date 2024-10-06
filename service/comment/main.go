package main

import (
	"diktok/config"
	pbcomment "diktok/grpc/comment"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/comment/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://comment", constant.ServiceName+".comment")
	// defer shutdown()
	database.InitMySQL()
	storage.CommentRedis = cache.InitRedis(config.System.Redis.CommentDB)
	etcd.InitETCD()
	ConnClose := rpc.InitRpcClient(etcd.GetEtcdClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServer(constant.CommentAddr, constant.CommentService, pbcomment.RegisterCommentServer, &CommentService{})
}
