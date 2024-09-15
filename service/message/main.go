package main

import (
	"diktok/config"
	pbmessage "diktok/grpc/message"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/database"
)

func main() {
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://message", constant.ServiceName+".message")
	// defer shutdown()
	database.InitMySQL()
	etcd.InitETCD()
	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClient(etcd.GetEtcdClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServer(constant.MessageAddr, constant.MessageService, pbmessage.RegisterMessageServer, &MessageService{})
}
