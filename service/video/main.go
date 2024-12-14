package main

import (
	"diktok/config"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/nacos"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/video/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	nacos.InitNacos()
	config.Init()
	util.InitZap()
	database.InitMySQL()
	// 初始化redis TODO收拢为一个redis
	storage.VideoRedis = cache.InitRedis(config.System.Redis.VideoDB)
	storage.UserRedis = cache.InitRedis(config.System.Redis.UserDB)
	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClientWithNacos(nacos.GetNamingClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServerWithNacos(constant.VideoAddr, constant.VideoService, pbvideo.RegisterVideoServer, &VideoService{})
}
