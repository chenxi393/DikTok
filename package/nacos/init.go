package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	local_constant "diktok/package/constant"
)

var (
	configClient config_client.IConfigClient
	namingClient naming_client.INamingClient
)

func InitNacos() {
	//create ServerConfig
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(local_constant.NacosIP, local_constant.NacosPort, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	clientConfig := constant.NewClientConfig(
		constant.WithNamespaceId(local_constant.NacosNameSpace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	// create naming client
	nc, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	namingClient = nc

	// create config client
	cc, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	configClient = cc
}

func GetConfigClient() config_client.IConfigClient {
	return configClient
}

func GetNamingClient() naming_client.INamingClient {
	return namingClient
}
