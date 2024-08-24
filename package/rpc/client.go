package rpc

import (
	pbcomment "diktok/grpc/comment"
	pbfavorite "diktok/grpc/favorite"
	pbmessage "diktok/grpc/message"
	pbrelation "diktok/grpc/relation"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"fmt"

	eclient "go.etcd.io/etcd/client/v3"
)

// 此为所有的RPC客户端 按需注册使用

var (
	UserClient     pbuser.UserClient
	VideoClient    pbvideo.VideoClient
	RelationClient pbrelation.RelationClient
	FavoriteClient pbfavorite.FavoriteClient
	MessageClient  pbmessage.MessageClient
	CommentClient  pbcomment.CommentClient
)

func InitRpcClient(etcdClient *eclient.Client) func() {
	funcSlice := make([]func() error, 0, 6)
	// RPC客户端 TODO 网关层最好设置一下RPC客户端 超时
	userConn := SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.UserService), etcdClient)
	UserClient = pbuser.NewUserClient(userConn)
	funcSlice = append(funcSlice, userConn.Close)

	videoConn := SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.VideoService), etcdClient)
	VideoClient = pbvideo.NewVideoClient(videoConn)
	funcSlice = append(funcSlice, videoConn.Close)

	favoriteConn := SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.FavoriteService), etcdClient)
	FavoriteClient = pbfavorite.NewFavoriteClient(favoriteConn)
	funcSlice = append(funcSlice, favoriteConn.Close)

	relationConn := SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.RalationService), etcdClient)
	RelationClient = pbrelation.NewRelationClient(relationConn)
	funcSlice = append(funcSlice, relationConn.Close)

	messageConn := SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.MessageService), etcdClient)
	MessageClient = pbmessage.NewMessageClient(messageConn)
	funcSlice = append(funcSlice, messageConn.Close)

	commentConn := SetupServiceConn(fmt.Sprintf("etcd:///%s", constant.CommentService), etcdClient)
	CommentClient = pbcomment.NewCommentClient(commentConn)
	funcSlice = append(funcSlice, commentConn.Close)

	return func() {
		for _, fun := range funcSlice {
			fun()
		}
	}
}