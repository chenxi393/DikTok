package rpc

import (
	"fmt"

	"github.com/chenxi393/nacos-grpc"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

// 弃用 使用nacos
func SetupServiceConn(serviceName string, etcdClient *eclient.Client) *grpc.ClientConn {
	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	// 然后在调用 grpc.Dial 方法创建连接代理 ClientConn 时，将其注入其中.
	// 类似于一个域名解析器
	etcdResolverBuilder, _ := eresolver.NewBuilder(etcdClient)
	// 开启用户服务的连接
	// Dial本身是不超时的（非阻塞的） 通过etcd拿到ip信息 etcd拿不到会阻塞
	// 这里的dial 是异步的 等到真正去调用 才会建立连接 而grpc 默认超时时间很长 需要手动设置
	conn, err := grpc.Dial(
		// 服务名称
		serviceName,
		// 注入 etcd resolverD
		grpc.WithResolvers(etcdResolverBuilder),
		// 声明使用的负载均衡策略为 roundrobin
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, roundrobin.Name)),
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

// TODO
func SetupServiceConnWithNacos(serviceName string, nacosClient naming_client.INamingClient) *grpc.ClientConn {
	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	// 然后在调用 grpc.Dial 方法创建连接代理 ClientConn 时，将其注入其中.
	// 类似于一个域名解析器
	resolverBuilder, _ := resolver.NewBuilder(nacosClient, "diktok")
	// 开启用户服务的连接
	// Dial本身是不超时的（非阻塞的） 通过etcd拿到ip信息 etcd拿不到会阻塞
	// 这里的dial 是异步的 等到真正去调用 才会建立连接 而grpc 默认超时时间很长 需要手动设置
	conn, err := grpc.NewClient(
		// 服务名称
		serviceName,
		// 注入 etcd resolverD
		grpc.WithResolvers(resolverBuilder),
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
