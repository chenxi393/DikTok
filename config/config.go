package config

import (
	"log"
	"strings"

	"diktok/package/constant"
	"diktok/package/nacos"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

type MySQL struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	UserName    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	MaxOpenConn int    `mapstructure:"maxOpenConn"`
	MaxIdleConn int    `mapstructure:"maxIdleConn"`
}

type HTTP struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	VideoAddress string `mapstructure:"videoAddress"`
}

type Redis struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	PoolSize   int    `mapstructure:"poolSize"`
	Password   string `mapstructure:"password"`
	UserDB     int    `mapstructure:"user_db"`
	VideoDB    int    `mapstructure:"video_db"`
	RelationDB int    `mapstructure:"relation_db"`
	FavoriteDB int    `mapstructure:"favorite_db"`
	CommentDB  int    `mapstructure:"comment_db"`
}

type RocketMQ struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type QiNiuCloud struct {
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	OssDomain string `mapstructure:"ossDomain"`
}

type OTEL struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
type SystemConfig struct {
	Mode          string     `mapstructure:"mode"`
	JwtSecret     string     `mapstructure:"jwtSecret"`
	GPTSecret     string     `mapstructure:"gptSecret"`
	Qiniu         QiNiuCloud `mapstructure:"qiniu"`
	HTTP          HTTP       `mapstructure:"http"`
	MysqlMaster   MySQL      `mapstructure:"mysqlMaster"`
	MysqlSlave    MySQL      `mapstructure:"mysqlSlave"`
	Redis         Redis      `mapstructure:"redis"`
	MQ            RocketMQ   `mapstructure:"rocketmq"`
	OtelColletcor OTEL       `mapstructure:"otel_collector"`
}

var System SystemConfig

// TODO 配置文件 各个服务应该分离
func Init() {
	//get config from nacos
	content, err := nacos.GetConfigClient().GetConfig(vo.ConfigParam{
		DataId: constant.NacosConfigId,
		Group:  constant.NacosGroupName,
	})
	if err != nil {
		log.Printf("[Init] get config from nacos failed: %s", err.Error())
	}
	if content == "" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config/")

		err = viper.ReadInConfig()
		if err != nil {
			log.Fatal("fatal error config file: ", err.Error())
		}

		err = viper.Unmarshal(&System)
		if err != nil {
			log.Fatal("fatal error unmarshal config: ", err.Error())
		}
		// 监视配置文件的变化 有变化就更改
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Println("配置文件被修改 重新载入全局变量")
			// 这里只有直接引用 config 才会热重载
			// 已经初始化好的 没啥用
			err = viper.Unmarshal(&System)
			if err != nil {
				log.Println("fatal error unmarshal config: ", err.Error())
			}
			log.Println(System.Qiniu.OssDomain)
		})
	} else { // 走配置中心
		viper.SetConfigType("yaml")
		err = viper.ReadConfig(strings.NewReader(content))
		if err != nil {
			log.Fatal("fatal error config file: ", err.Error())
		}
		err = viper.Unmarshal(&System)
		if err != nil {
			log.Fatal("fatal error unmarshal config: ", err.Error())
		}

		//Listen config change,key=dataId+group+namespaceId.
		err = nacos.GetConfigClient().ListenConfig(vo.ConfigParam{
			DataId: "config.yaml",
			Group:  "diktok",
			OnChange: func(namespace, group, dataId, data string) {
				log.Println("配置文件被修改 重新载入全局变量")
				err = viper.ReadConfig(strings.NewReader(data))
				if err != nil {
					log.Fatal("fatal error config file: ", err.Error())
				}
				err = viper.Unmarshal(&System)
				if err != nil {
					log.Println("fatal error unmarshal config: ", err.Error())
				}
				log.Println(System.Qiniu.OssDomain)
			},
		})
	}

	log.Println("[Init] viper读取配置文件成功")
}
