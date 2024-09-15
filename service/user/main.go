package main

import (
	"diktok/config"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/user/logic"
	"diktok/service/user/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	// 初始化配置文件
	config.Init()
	// 初始化日志打印
	util.InitZap()
	// shutdown := otel.Init("rpc://user", constant.ServiceName+".user")
	// defer shutdown()
	database.InitMySQL()
	storage.UserRedis = cache.InitRedis(config.System.Redis.UserDB)
	logic.RegisterChatGPT()
	// 初始化ETCD 作为服务发现与注册中心
	etcd.InitETCD()
	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClient(etcd.GetEtcdClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServer(constant.UserAddr, constant.UserService, pbuser.RegisterUserServer, &UserService{})
}
