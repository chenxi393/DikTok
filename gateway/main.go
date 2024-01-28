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
	"douyin/package/otel"
	"douyin/package/rpc"
	"douyin/package/util"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	resolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.uber.org/zap"
)

func main() {
	config.Init()
	util.InitZap()
	shutdown := otel.Init("http://newclip.cn", constant.ServiceName+".gateway")
	defer shutdown()
	// 创建 etcd 客户端 先与服务建立连接
	etcdClient, err := clientv3.NewFromURL(constant.MyEtcdURL)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	// 然后在调用 grpc.Dial 方法创建连接代理 ClientConn 时，将其注入其中.
	// 类似于一个域名解析器
	etcdResolverBuilder, err := resolver.NewBuilder(etcdClient)
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
	userConn := rpc.SetupServiceConn(userTarget, etcdResolverBuilder)
	handler.UserClient = pbuser.NewUserClient(userConn)
	defer userConn.Close()

	videoConn := rpc.SetupServiceConn(videoTarget, etcdResolverBuilder)
	handler.VideoClient = pbvideo.NewVideoClient(videoConn)
	defer userConn.Close()

	relationConn := rpc.SetupServiceConn(relationTarget, etcdResolverBuilder)
	handler.RelationClient = pbrelation.NewRelationClient(relationConn)
	defer userConn.Close()

	favoriteConn := rpc.SetupServiceConn(favoriteTarget, etcdResolverBuilder)
	handler.FavoriteClient = pbfavorite.NewFavoriteClient(favoriteConn)
	defer userConn.Close()

	messageConn := rpc.SetupServiceConn(messageTarget, etcdResolverBuilder)
	handler.MessageClinet = pbmessage.NewMessageClient(messageConn)
	defer userConn.Close()

	commentConn := rpc.SetupServiceConn(commentTarget, etcdResolverBuilder)
	handler.CommentClient = pbcomment.NewCommentClient(commentConn)
	defer userConn.Close()

	// 初始化http框架 并listen
	startFiber()
}
