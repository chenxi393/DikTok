package rpc

import (
	"context"
	"log"
	"net"
	"strconv"
	"strings"

	"diktok/package/etcd"
	"diktok/package/nacos"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func InitServer[T interface{}](addr, serviceName string, registrFunc func(s grpc.ServiceRegistrar, srv T), svrHandler interface{}) {
	if etcd.GetEtcdClient() == nil {
		panic("etcd client is nil")
	}
	// RPC服务端D
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(30*1024*1024), grpc.StatsHandler(otelgrpc.NewServerHandler()))
	// 如果一个结构体实现了interface 则该结构体对象 和对象指针 均可以类型断言为interface
	registrFunc(s, svrHandler.(T))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册 grpc 服务节点到 etcd 中
	go RegisterEndPointToEtcd(ctx, etcd.GetEtcdClient(), addr, serviceName)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func InitServerWithNacos[T interface{}](addr, serviceName string, registrFunc func(s grpc.ServiceRegistrar, srv T), svrHandler interface{}) {
	if nacos.GetNamingClient() == nil {
		panic("nacos client is nil")
	}
	var (
		lis net.Listener
		err error
	)
	for i := 0; i < 10; i++ {
		lis, err = net.Listen("tcp", addr)
		if err != nil {
			log.Printf("failed to listen: %v", err)
			port := strings.Split(addr, ":")[1]
			portInt, _ := strconv.Atoi(port)
			addr = strings.Replace(addr, port, strconv.Itoa(portInt+1), 1)
		} else {
			break
		}
	}
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(30*1024*1024), grpc.StatsHandler(otelgrpc.NewServerHandler()))
	// 如果一个结构体实现了interface 则该结构体对象 和对象指针 均可以类型断言为interface
	registrFunc(s, svrHandler.(T))

	go registerToNacos(nacos.GetNamingClient(), addr, serviceName)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
