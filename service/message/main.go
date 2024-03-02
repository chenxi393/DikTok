package main

import (
	"context"
	"douyin/config"
	pbmessage "douyin/grpc/message"
	pbrelation "douyin/grpc/relation"
	"douyin/package/constant"
	"douyin/package/otel"
	"douyin/package/rpc"
	"douyin/package/util"
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
)

func main() {
	config.Init()
	util.InitZap()
	shutdown := otel.Init("rpc://message", constant.ServiceName+".message")
	defer shutdown()
	close := InitMongoDB()
	defer close()

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

	lis, err := net.Listen("tcp", constant.MessageAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbmessage.RegisterMessageServer(s, &MessageService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册grpc到etcd节点中
	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, constant.MessageAddr, constant.MessageService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
