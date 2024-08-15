package main

import (
	"context"
	"log"
	"net"

	"diktok/config"
	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/package/etcd"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/cache"
	"diktok/storage/database"

	"github.com/go-redis/redis"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var (
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
	etcd.InitETCD()

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
	go rpc.RegisterEndPointToEtcd(ctx, etcd.GetEtcdClient(), constant.FavoriteAddr, constant.FavoriteService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
