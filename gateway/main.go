package main

import (
	"fmt"

	"time"

	"diktok/config"
	"diktok/gateway/handler"
	pbcomment "diktok/grpc/comment"
	pbfavorite "diktok/grpc/favorite"
	pbmessage "diktok/grpc/message"
	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	config.Init()
	util.InitZap()
	// shutdown := otel.Init("http://newclip.cn", constant.ServiceName+".gateway")
	// defer shutdown()

	// 连接ETCD
	// TODO 非常奇怪 连接服务器的etcd 会超时  linux 系统 https 走socks 代理的问题
	// etcd老版本使用的  grpc版本 默认不走代理
	// etcd 客户端与grpc服务端 通信使用的grpc协议 可以看下etcd的源码
	// 因为grpc建立连接非阻塞调用 所以在put之类的操作的时候会阻塞
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:            []string{config.System.EtcdURL},
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
		DialTimeout:          5 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
	})
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// RPC客户端 TODO 网关层最好设置一下RPC客户端 超时
	userConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.UserService), etcdClient)
	handler.UserClient = pbuser.NewUserClient(userConn)
	defer userConn.Close()

	videoConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.VideoService), etcdClient)
	handler.VideoClient = pbvideo.NewVideoClient(videoConn)
	defer videoConn.Close()

	relationConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.RalationService), etcdClient)
	handler.RelationClient = pbrelation.NewRelationClient(relationConn)
	defer relationConn.Close()

	favoriteConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.FavoriteService), etcdClient)
	handler.FavoriteClient = pbfavorite.NewFavoriteClient(favoriteConn)
	defer favoriteConn.Close()

	messageConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.MessageService), etcdClient)
	handler.MessageClinet = pbmessage.NewMessageClient(messageConn)
	defer messageConn.Close()

	commentConn := rpc.SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.CommentService), etcdClient)
	handler.CommentClient = pbcomment.NewCommentClient(commentConn)
	defer commentConn.Close()

	// 开启HTTP框架
	startFiber()
}
