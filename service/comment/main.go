package main

import (
	"context"
	"douyin/config"
	pbcomment "douyin/grpc/comment"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"douyin/package/rpc"
	"douyin/package/util"
	"douyin/storage/cache"
	"douyin/storage/database"
	"douyin/storage/mq"
	"fmt"
	"log"
	"net"

	eclient "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// 需要RPC调用的客户端
	userClient pbuser.UserClient
	// comment模块运行在 8060-8069
	addr = "comment:8060"
)

func main() {
	// TODO 配置文件实际上也应该分离
	config.Init()
	// TODO 日志也应该考虑合并
	util.InitZap()
	database.InitMySQL()
	cache.InitRedis()
	mq.InitComment()

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
	userTarget := fmt.Sprintf("etcd:///%s", constant.UserService)
	userConn := rpc.SetupServiceConn(userTarget, etcdResolverBuilder)
	defer userConn.Close()
	userClient = pbuser.NewUserClient(userConn)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pbcomment.RegisterCommentServer(s, &CommentService{})

	// TODO 这一块context 目前还没没有理解是干嘛的
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册grpc到etcd节点中
	// 注册 grpc 服务节点到 etcd 中
	go rpc.RegisterEndPointToEtcd(ctx, addr, constant.CommentService)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
