package main

import (
	"context"
	"log"
	"net"

	"diktok/config"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/video/storage"
	"diktok/storage/cache"
	"diktok/storage/database"

	eclient "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
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
	defer rpc.InitRpcClient(etcd.GetEtcdClient())
	// 初始化rpc 服务端
	InitServer(etcd.GetEtcdClient())
}

func InitServer(etcdClient *eclient.Client) {
	// RPC服务端
	lis, err := net.Listen("tcp", constant.VideoAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(30 * 1024 * 1024))
	pbvideo.RegisterVideoServer(s, &VideoService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, etcdClient, constant.VideoAddr, constant.VideoService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
