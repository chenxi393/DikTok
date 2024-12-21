package constant

import (
	"os"
	"strings"
)

func init() {
	// TODO 适配本机和docker 是不是有更好的办法
	if os.Getenv("RUN_ENV") != "docker" {
		VideoAddr = getLocalAddr(VideoAddr)
		UserAddr = getLocalAddr(UserAddr)
		RelationAddr = getLocalAddr(RelationAddr)
		MessageAddr = getLocalAddr(MessageAddr)
		FavoriteAddr = getLocalAddr(FavoriteAddr)
		CommentAddr = getLocalAddr(CommentAddr)
		NacosIP = "127.0.0.1"
		NacosNameSpace = "diktok-offline"
	}
}

// grpc
var (
	UserService     = "diktok/user"
	VideoService    = "diktok/video"
	RalationService = "diktok/relation"
	CommentService  = "diktok/comment"
	MessageService  = "diktok/message"
	FavoriteService = "diktok/favorite"

	VideoAddr    = "video:8010" //docker 内使用video:8010 本地使用127.0.0.1 否则会阻塞
	UserAddr     = "user:8020"
	RelationAddr = "relation:8030"
	MessageAddr  = "message:8040"
	FavoriteAddr = "favorite:8050"
	CommentAddr  = "comment:8060"

	NacosNameSpace        = "diktok-online"
	NacosGroupName        = "diktok"
	NacosConfigId         = "config_common"
	NacosIP               = "nacos"
	NacosPort      uint64 = 8848
)

func getLocalAddr(addr string) string {
	e := strings.Split(addr, ":")
	return "127.0.0.1:" + e[1]
}
