package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"diktok/config"
	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/cache"
	"diktok/storage/database"

	"github.com/go-redis/redis"
	eclient "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// 需要RPC调用的客户端
	relationClient pbrelation.RelationClient
	userRedis      *redis.Client
)

func main() {
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("rpc://user", constant.ServiceName+".user")
	// defer shutdown()
	database.InitMySQL()
	userRedis = cache.InitRedis(config.System.Redis.UserDB)
	registerChatGPT()

	// 连接ETCD
	etcdClient, err := eclient.NewFromURL(config.System.EtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// RPC客户端
	relationConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.RalationService), etcdClient)
	defer relationConn.Close()
	relationClient = pbrelation.NewRelationClient(relationConn)

	// RPC服务端
	lis, err := net.Listen("tcp", constant.UserAddr)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	// 添加 grpc otel 自动检测
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbuser.RegisterUserServer(s, &UserService{})

	// TODO 这一块context 目前还没没有理解是干嘛的
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, etcdClient, constant.UserAddr, constant.UserService)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
