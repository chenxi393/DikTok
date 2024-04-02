package rpc

import (
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func SetupServiceConn(serviceName string, resolver resolver.Builder) *grpc.ClientConn {
	// 开启用户服务的连接
	// Dial本身是不超时的（非阻塞的） 也可以通过etcd拿到ip信息 只不过ip没有部署 会阻塞
	// 这里的dial 是异步的 等到真正去调用 才会建立连接 而grpc 默认超时时间很长 需要手动设置
	conn, err := grpc.Dial(
		// 服务名称
		serviceName,
		// 注入 etcd resolverD
		grpc.WithResolvers(resolver),
		// 声明使用的负载均衡策略为 roundrobin
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc自动检测
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		// 不推荐使用 阻塞
		// grpc.WithBlock()
	)
	if err != nil {
		zap.L().Sugar().Fatalf("did not connect %s : %v\n", serviceName, err)
	}
	return conn
}
