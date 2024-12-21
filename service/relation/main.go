package main

import (
	"diktok/config"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/package/nacos"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/relation/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	nacos.InitNacos()
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://relation", constant.ServiceName+".relation")
	// defer shutdown()
	database.InitMySQL()
	storage.RelationRedis = cache.InitRedis(config.System.Redis.RelationDB)
	storage.UserRedis = cache.InitRedis(config.System.Redis.UserDB)

	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClientWithNacos(nacos.GetNamingClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServerWithNacos(constant.RelationAddr, constant.RalationService, pbrelation.RegisterRelationServer, &RelationService{})
}
