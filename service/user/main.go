package main

import (
	"context"
	"douyin/config"
	"douyin/database"
	pbuser "douyin/grpc/user"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/llm"
	"douyin/package/mq"
	"douyin/package/util"
	"fmt"
	"log"
	"net"
	"time"

	eclient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

const (
	// etcd 端口
	MyEtcdURL = "http://localhost:2379"

	addr = "127.0.0.1:6668"
)

func main() {
	// TODO 配置文件实际上也应该分离
	config.Init()
	// TODO 日志也应该考虑合并
	util.InitZap()
	database.InitMySQL()
	cache.InitRedis()
	mq.InitMQ()
	llm.RegisterChatGPT()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pbuser.RegisterUserServer(s, &UserService{})

	// TODO 这一块context 目前还没没有理解是干嘛的
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 注册grpc到etcd节点中
	// 注册 grpc 服务节点到 etcd 中
	go registerEndPointToEtcd(ctx, addr)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func registerEndPointToEtcd(ctx context.Context, addr string) {
	// 创建 etcd 客户端
	etcdClient, _ := eclient.NewFromURL(MyEtcdURL)
	// 创建 etcd 服务端节点管理模块 etcdManager
	etcdManager, _ := endpoints.NewManager(etcdClient, constant.UserService)

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	var ttl int64 = 10
	lease, _ := etcdClient.Grant(ctx, ttl)

	// 添加注册节点到 etcd 中，并且携带上租约 id
	_ = etcdManager.AddEndpoint(ctx, fmt.Sprintf("%s/%s", constant.UserService, addr),
		endpoints.Endpoint{Addr: addr}, eclient.WithLease(lease.ID))

	// 每隔 5 s进行一次延续租约的动作
	for {
		select {
		case <-time.After(5 * time.Second):
			// 续约操作
			resp, _ := etcdClient.KeepAliveOnce(ctx, lease.ID)
			log.Printf("keep alive resp: %+v\n", resp)
		case <-ctx.Done():
			return
		}
	}
}
