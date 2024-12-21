package main

import (
	"diktok/config"
	pbmessage "diktok/grpc/message"
	"diktok/package/constant"
	"diktok/package/nacos"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/database"
)

func main() {
	nacos.InitNacos()
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://message", constant.ServiceName+".message")
	// defer shutdown()
	database.InitMySQL()
	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClientWithNacos(nacos.GetNamingClient())
	defer ConnClose()
	// 初始化rpc 服务端
	rpc.InitServerWithNacos(constant.MessageAddr, constant.MessageService, pbmessage.RegisterMessageServer, &MessageService{})
}
