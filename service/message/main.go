package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"diktok/config"
	pbmessage "diktok/grpc/message"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/package/otel"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/database"

	eclient "go.etcd.io/etcd/client/v3"
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
	database.InitMySQL()

	// 连接ETCD
	etcdClient, err := eclient.NewFromURL(config.System.EtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// RPC客户端
	relationConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.RalationService), etcdClient)
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

	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, etcdClient, constant.MessageAddr, constant.MessageService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
