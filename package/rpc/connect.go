package rpc

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func SetupServiceConn(serviceName string, resolver resolver.Builder) *grpc.ClientConn {
	// 开启用户服务的连接
	conn, err := grpc.Dial(
		// 服务名称
		serviceName,
		// 注入 etcd resolverD
		grpc.WithResolvers(resolver),
		// 声明使用的负载均衡策略为 roundrobin
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Sugar().Fatalf("did not connect %s : %v\n", serviceName, err)
	}
	return conn
}
