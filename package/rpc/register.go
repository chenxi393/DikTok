package rpc

import (
	"context"
	"diktok/package/constant"
	"log"
	"strconv"
	"time"

	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	eclient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

func RegisterEndPointToEtcd(ctx context.Context, etcdClient *eclient.Client, addr string, name string) {
	// 创建 etcd 服务端节点管理模块 etcdManager
	etcdManager, _ := endpoints.NewManager(etcdClient, name)

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	var ttl int64 = 10
	lease, _ := etcdClient.Grant(ctx, ttl)

	// 添加注册节点到 etcd 中，并且携带上租约 id
	err := etcdManager.AddEndpoint(ctx, name+"/"+addr, endpoints.Endpoint{Addr: addr}, eclient.WithLease(lease.ID))
	if err != nil {
		log.Fatalf("add endpoint err: %v", err)
	}

	// 每隔 5 s进行一次延续租约的动作
	for {
		select {
		case <-time.After(5 * time.Second):
			// 续约操作
			_, err := etcdClient.KeepAliveOnce(ctx, lease.ID)
			//log.Printf("keep alive resp: %+v\n", resp)
			if err != nil {
				// 直接 fatal 然后etcd 重启就行
				log.Fatalf("keep alive err: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func registerToNacos(client naming_client.INamingClient, addr string, name string) {
	addrSlice := strings.Split(addr, ":")
	if len(addrSlice) != 2 {
		log.Fatalf("addr is invalid: %v", addr)
	}
	port, err := strconv.ParseUint(addrSlice[1], 10, 64)
	if err != nil {
		log.Fatalf("addr is invalid: %v", addr)
	}
	registerServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          addrSlice[0],
		Port:        port,
		ServiceName: name,
		GroupName:   constant.NacosGroupName,
		ClusterName: "default",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true, // 临时节点
		Metadata:    map[string]string{"idc": "hongkong"},
	})
}

func registerServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	success, err := client.RegisterInstance(param)
	if !success || err != nil {
		panic("RegisterServiceInstance failed!" + err.Error())
	}
	log.Printf("RegisterServiceInstance,param:%+v,result:%+v", param, success)
}
