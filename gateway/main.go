package main

import (
	"douyin/config"
	"douyin/gateway/handler"
	pbcomment "douyin/grpc/comment"
	pbfavorite "douyin/grpc/favorite"
	pbmessage "douyin/grpc/message"
	pbrelation "douyin/grpc/relation"
	pbuser "douyin/grpc/user"
	pbvideo "douyin/grpc/video"
	"douyin/package/constant"
	"douyin/package/util"
	"fmt"

	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func main() {
	// TODO 配置文件实际上也应该分离
	config.Init()
	// TODO 日志也应该考虑合并
	util.InitZap()
	// 先与服务建立连接
	connectService()
	// 再初始化http框架 并listen
	startFiber()
}

func connectService() {
	// 创建 etcd 客户端
	etcdClient, err := eclient.NewFromURL(constant.MyEtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	// 然后在调用 grpc.Dial 方法创建连接代理 ClientConn 时，将其注入其中.
	// 类似于一个域名解析器
	etcdResolverBuilder, err := eresolver.NewBuilder(etcdClient)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	// 拼接服务名称，需要固定义 etcd:/// 作为前缀
	userTarget := fmt.Sprintf("etcd:///%s", constant.UserService)
	videoTarget := fmt.Sprintf("etcd:///%s", constant.VideoService)
	relationTarget := fmt.Sprintf("etcd:///%s", constant.RalationService)
	favoriteTarget := fmt.Sprintf("etcd:///%s", constant.FavoriteService)
	commentTarget := fmt.Sprintf("etcd:///%s", constant.CommentService)
	messageTarget := fmt.Sprintf("etcd:///%s", constant.MessageService)

	// 开启用户服务的连接 并且defer关闭函数
	userConn := setupServiceConn(userTarget, etcdResolverBuilder)
	handler.UserClient = pbuser.NewUserClient(userConn)
	// FIXME 这里不能用close
	// defer userConn.Close()

	videoConn := setupServiceConn(videoTarget, etcdResolverBuilder)
	handler.VideoClient = pbvideo.NewVideoClient(videoConn)
	// defer userConn.Close()

	relationConn := setupServiceConn(relationTarget, etcdResolverBuilder)
	handler.RelationClient = pbrelation.NewRelationClient(relationConn)
	// defer userConn.Close()

	favoriteConn := setupServiceConn(favoriteTarget, etcdResolverBuilder)
	handler.FavoriteClient = pbfavorite.NewFavoriteClient(favoriteConn)
	// defer userConn.Close()

	messageConn := setupServiceConn(messageTarget, etcdResolverBuilder)
	handler.MessageClinet = pbmessage.NewMessageClient(messageConn)
	// defer userConn.Close()

	commentConn := setupServiceConn(commentTarget, etcdResolverBuilder)
	handler.CommentClient = pbcomment.NewCommentClient(commentConn)
	// defer userConn.Close()
}

func setupServiceConn(serviceName string, resolver resolver.Builder) *grpc.ClientConn {
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
