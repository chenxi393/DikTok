package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"diktok/config"
	pbfavorite "diktok/grpc/favorite"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/cache"
	"diktok/storage/database"

	"github.com/go-redis/redis"
	eclient "go.etcd.io/etcd/client/v3"
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
	// shutdown := otel.Init("rpc://favorite", constant.ServiceName+".favorite")
	// defer shutdown()
	database.InitMySQL()
	favoriteRedis = cache.InitRedis(config.System.Redis.FavoriteDB)
	videoRedis = cache.InitRedis(config.System.Redis.VideoDB)
	userRedis = cache.InitRedis(config.System.Redis.UserDB)

	// 连接ETCD
	etcdClient, err := eclient.NewFromURL(config.System.EtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// RPC客户端
	videoConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.VideoService), etcdClient)
	defer videoConn.Close()
	videoClient = pbvideo.NewVideoClient(videoConn)

	// RPC服务端
	lis, err := net.Listen("tcp", constant.FavoriteAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pbfavorite.RegisterFavoriteServer(s, &FavoriteService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, etcdClient, constant.FavoriteAddr, constant.FavoriteService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
