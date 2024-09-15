package main

import (
	"diktok/config"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/video/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	// 初始化配置文件
	config.Init()
	// 初始化日志打印
	util.InitZap()
	// 初始化DB
	database.InitMySQL()
	// 初始化redis TODO收拢为一个redis
	storage.VideoRedis = cache.InitRedis(config.System.Redis.VideoDB)
	storage.UserRedis = cache.InitRedis(config.System.Redis.UserDB)
	// 初始化ETCD 作为服务发现与注册中心
	etcd.InitETCD()
	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClient(etcd.GetEtcdClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServer(constant.VideoAddr, constant.VideoService, pbvideo.RegisterVideoServer, &VideoService{})
}
