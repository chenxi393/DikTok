package main

import (
	"context"
	"douyin/config"
	pbmessage "douyin/grpc/message"
	pbrelation "douyin/grpc/relation"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"douyin/package/otel"
	"douyin/package/rpc"
	"douyin/package/util"
	"douyin/storage/cache"
	"douyin/storage/database"
	"douyin/storage/mq"
	"fmt"
	"log"
	"net"

	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// 需要RPC调用的客户端
	userClient    pbuser.UserClient
	messageClient pbmessage.MessageClient
	// relation模块运行在 8030-8039
	addr = "127.0.0.1:8030"
)

func main() {
	config.Init()
	util.InitZap()
	shutdown := otel.Init("rpc://relation", constant.ServiceName+".relation")
	defer shutdown()
	database.InitMySQL()
	cache.InitRedis()
	mq.InitRelation()
	go mq.FollowConsume()

	// 连接到依赖的服务
	etcdClient, err := eclient.NewFromURL(constant.MyEtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	etcdResolverBuilder, err := eresolver.NewBuilder(etcdClient)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	// 开启用户服务的连接 并且defer关闭函数
	userTarget := fmt.Sprintf("etcd:///%s", constant.UserService)
	userConn := rpc.SetupServiceConn(userTarget, etcdResolverBuilder)
	defer userConn.Close()
	userClient = pbuser.NewUserClient(userConn)

	messageTarget := fmt.Sprintf("etcd:///%s", constant.MessageService)
	messageConn := rpc.SetupServiceConn(messageTarget, etcdResolverBuilder)
	defer messageConn.Close()
	messageClient = pbmessage.NewMessageClient(messageConn)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbrelation.RegisterRelationServer(s, &RelationService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册grpc到etcd节点中
	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, addr, constant.RalationService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}