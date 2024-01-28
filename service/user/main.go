package main

import (
	"context"
	"douyin/config"
	pbrelation "douyin/grpc/relation"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"douyin/package/otel"
	"douyin/package/rpc"
	"douyin/package/util"
	"douyin/storage/cache"
	"douyin/storage/database"
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
	relationClient pbrelation.RelationClient
	// 用户模块运行在 8020-8029  FIXME 这里用 user:8020 直接报错。。 user会直接解析出端口
	addr = "127.0.0.1:8020"
)

func main() {
	config.Init()
	util.InitZap()
	shutdown := otel.Init("rpc://user", constant.ServiceName+".user")
	defer shutdown()
	database.InitMySQL()
	cache.InitRedis()
	registerChatGPT()

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
	relationTarget := fmt.Sprintf("etcd:///%s", constant.RalationService)
	relationConn := rpc.SetupServiceConn(relationTarget, etcdResolverBuilder)
	defer relationConn.Close()
	relationClient = pbrelation.NewRelationClient(relationConn)
	// 注册自己的服务
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	// 添加 grpc otel 自动检测
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbuser.RegisterUserServer(s, &UserService{})
	// TODO 这一块context 目前还没没有理解是干嘛的
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册grpc到etcd节点中
	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, addr, constant.UserService)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
