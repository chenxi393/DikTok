package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"diktok/config"
	pbfavorite "diktok/grpc/favorite"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/storage/cache"
	"diktok/storage/database"

	"github.com/go-redis/redis"
	eclient "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	userClient     pbuser.UserClient
	favoriteClient pbfavorite.FavoriteClient

	videoRedis, userRedis *redis.Client
)

func main() {
	config.Init()
	util.InitZap()
	database.InitMySQL()
	videoRedis = cache.InitRedis(config.System.Redis.VideoDB)
	userRedis = cache.InitRedis(config.System.Redis.UserDB)

	// 连接ETCD
	etcdClient, err := eclient.NewFromURL(config.System.EtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// RPC客户端
	userConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.UserService), etcdClient)
	defer userConn.Close()
	userClient = pbuser.NewUserClient(userConn)

	favoriteConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.FavoriteService), etcdClient)
	favoriteClient = pbfavorite.NewFavoriteClient(favoriteConn)
	defer favoriteConn.Close()

	// RPC服务端
	lis, err := net.Listen("tcp", constant.VideoAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(30 * 1024 * 1024))
	pbvideo.RegisterVideoServer(s, &VideoService{})

	// TODO 这一块context 目前还没没有理解是干嘛的
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, etcdClient, constant.VideoAddr, constant.VideoService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
