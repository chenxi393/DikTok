package main

import (
	"context"
	"douyin/config"
	pbfavorite "douyin/grpc/favorite"
	pbvideo "douyin/grpc/video"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/database"
	"douyin/package/otel"
	"douyin/package/rpc"
	"douyin/package/util"
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis"
	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// 需要RPC调用的客户端
	videoClient pbvideo.VideoClient

	favoriteRedis, videoRedis, userRedis *redis.Client
)

func main() {
	config.Init()
	util.InitZap()
	shutdown := otel.Init("rpc://favorite", constant.ServiceName+".favorite")
	defer shutdown()
	database.InitMySQL()
	favoriteRedis = cache.InitRedis(config.System.Redis.FavoriteDB)
	videoRedis = cache.InitRedis(config.System.Redis.VideoDB)
	userRedis = cache.InitRedis(config.System.Redis.UserDB)

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
	videoTarget := fmt.Sprintf("etcd:///%s", constant.VideoService)
	videoConn := rpc.SetupServiceConn(videoTarget, etcdResolverBuilder)
	defer videoConn.Close()
	videoClient = pbvideo.NewVideoClient(videoConn)

	lis, err := net.Listen("tcp", constant.FavoriteAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbfavorite.RegisterFavoriteServer(s, &FavoriteService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册grpc到etcd节点中
	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, constant.FavoriteAddr, constant.FavoriteService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
