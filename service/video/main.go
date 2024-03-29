package main

import (
	"context"
	"douyin/config"
	pbfavorite "douyin/grpc/favorite"
	pbuser "douyin/grpc/user"
	pbvideo "douyin/grpc/video"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/database"
	"douyin/package/rpc"
	"douyin/package/util"
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis"
	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
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

	// 创建 etcd 客户端 连接到
	etcdClient, err := eclient.NewFromURL(constant.MyEtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	etcdResolverBuilder, err := eresolver.NewBuilder(etcdClient)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	// 拼接服务名称，需要固定义 etcd:/// 作为前缀
	userTarget := fmt.Sprintf("etcd:///%s", constant.UserService)
	favoriteTarget := fmt.Sprintf("etcd:///%s", constant.FavoriteService)

	// 开启用户服务的连接 并且defer关闭函数
	userConn := rpc.SetupServiceConn(userTarget, etcdResolverBuilder)
	userClient = pbuser.NewUserClient(userConn)
	defer userConn.Close()

	// 开启用户服务的连接 并且defer关闭函数
	favoriteConn := rpc.SetupServiceConn(favoriteTarget, etcdResolverBuilder)
	favoriteClient = pbfavorite.NewFavoriteClient(favoriteConn)
	defer favoriteConn.Close()

	lis, err := net.Listen("tcp", constant.VideoAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(30 * 1024 * 1024))
	pbvideo.RegisterVideoServer(s, &VideoService{})

	// TODO 这一块context 目前还没没有理解是干嘛的
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册grpc到etcd节点中
	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, constant.VideoAddr, constant.VideoService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
