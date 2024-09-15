package main

import (
	"diktok/config"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/relation/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
)

func main() {
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://relation", constant.ServiceName+".relation")
	// defer shutdown()
	database.InitMySQL()
	storage.RelationRedis = cache.InitRedis(config.System.Redis.RelationDB)
	storage.UserRedis = cache.InitRedis(config.System.Redis.UserDB)
	etcd.InitETCD()
	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClient(etcd.GetEtcdClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServer(constant.RelationAddr, constant.RalationService, pbrelation.RegisterRelationServer, &RelationService{})
}
