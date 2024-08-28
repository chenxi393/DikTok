package main

import (
	"diktok/config"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
)

func main() {
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("http://newclip.cn", constant.ServiceName+".gateway")
	// defer shutdown()
	etcd.InitETCD()
	ConnClose := rpc.InitRpcClient(etcd.GetEtcdClient())
	defer ConnClose()
	// 开启HTTP框架
	startFiber()
}
