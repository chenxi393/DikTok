package main

import (
	"diktok/config"
	"diktok/package/nacos"
	"diktok/package/rpc"
	"diktok/package/util"
)

func main() {
	nacos.InitNacos()
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("http://newclip.cn", constant.ServiceName+".gateway")
	// defer shutdown()
	ConnClose := rpc.InitRpcClientWithNacos(nacos.GetNamingClient())
	defer ConnClose()
	// 开启HTTP框架
	startFiber()
}
