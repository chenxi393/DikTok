package rpc

import (
	"context"
	"douyin/package/constant"
	"fmt"
	"log"
	"time"

	eclient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

func RegisterEndPointToEtcd(ctx context.Context, addr string, name string) {
	// 创建 etcd 客户端
	etcdClient, _ := eclient.NewFromURL(constant.MyEtcdURL)
	// 创建 etcd 服务端节点管理模块 etcdManager
	etcdManager, _ := endpoints.NewManager(etcdClient, name)

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	var ttl int64 = 10
	lease, _ := etcdClient.Grant(ctx, ttl)

	// 添加注册节点到 etcd 中，并且携带上租约 id
	_ = etcdManager.AddEndpoint(ctx, fmt.Sprintf("%s/%s", name, addr),
		endpoints.Endpoint{Addr: addr}, eclient.WithLease(lease.ID))

	// 每隔 5 s进行一次延续租约的动作
	for {
		select {
		case <-time.After(5 * time.Second):
			// 续约操作
			_, err := etcdClient.KeepAliveOnce(ctx, lease.ID)
			//log.Printf("keep alive resp: %+v\n", resp)
			if err != nil {
				log.Printf("keep alive err: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
