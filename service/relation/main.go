package main

import (
	"context"
	"log"
	"net"

	"diktok/config"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/cache"
	"diktok/storage/database"

	"github.com/go-redis/redis"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var (
	relationRedis, userRedis *redis.Client
)

func main() {
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://relation", constant.ServiceName+".relation")
	// defer shutdown()
	database.InitMySQL()
	relationRedis = cache.InitRedis(config.System.Redis.RelationDB)
	userRedis = cache.InitRedis(config.System.Redis.UserDB)
	etcd.InitETCD()
	// 初始化rpc 客户端
	ConnClose := rpc.InitRpcClient(etcd.GetEtcdClient())
	defer ConnClose()
	// RPC服务端
	lis, err := net.Listen("tcp", constant.RelationAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbrelation.RegisterRelationServer(s, &RelationService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, etcd.GetEtcdClient(), constant.RelationAddr, constant.RalationService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
