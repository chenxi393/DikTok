package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"diktok/config"
	pbmessage "diktok/grpc/message"
	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/otel"
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
	userClient    pbuser.UserClient
	messageClient pbmessage.MessageClient

	relationRedis, userRedis *redis.Client
)

func main() {
	config.Init()
	util.InitZap()
	shutdown := otel.Init("rpc://relation", constant.ServiceName+".relation")
	defer shutdown()
	database.InitMySQL()
	relationRedis = cache.InitRedis(config.System.Redis.RelationDB)
	userRedis = cache.InitRedis(config.System.Redis.UserDB)

	// 连接ETCD
	etcdClient, err := eclient.NewFromURL(config.System.EtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// RPC客户端
	userConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.UserService), etcdClient)
	defer userConn.Close()
	userClient = pbuser.NewUserClient(userConn)

	messageConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.MessageService), etcdClient)
	defer messageConn.Close()
	messageClient = pbmessage.NewMessageClient(messageConn)

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
	go rpc.RegisterEndPointToEtcd(ctx, etcdClient, constant.RelationAddr, constant.RalationService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
