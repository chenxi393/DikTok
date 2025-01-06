package etcd

// import (
// 	"diktok/config"

// 	eclient "go.etcd.io/etcd/client/v3"
// 	"go.uber.org/zap"
// )

// var etcdClient *eclient.Client

// // 已弃用 使用nacos作为注册中心
// func InitETCD() {
// 	// 连接ETCD
// 	// TODO 新版本SDK连接服务器的etcd 会超时  linux 系统 https 走socks 代理的问题
// 	// etcd老版本SDK使用的 grpc版本 默认不走代理
// 	var err error
// 	etcdClient, err = eclient.NewFromURL(config.System.EtcdURL)
// 	if err != nil {
// 		zap.L().Fatal(err.Error())
// 	}
// }

// func GetEtcdClient() *eclient.Client {
// 	return etcdClient
// }
