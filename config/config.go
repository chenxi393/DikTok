package config

import (
	"log"
	"os"
	"strings"

	"diktok/package/constant"

	"github.com/fsnotify/fsnotify"
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

type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type MongoDB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
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
	MQ            RabbitMQ   `mapstructure:"rabbitmq"`
	OtelColletcor OTEL       `mapstructure:"otel_collector"`
	Mongo         MongoDB    `mapstructure:"mongo"`
	EtcdURL       string     `mapstructure:"etcd_address"`
}

var System SystemConfig

// TODO 配置文件 各个服务应该分离
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
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

	// TODO 适配本机和docker 是不是有更好的办法
	if os.Getenv("RUN_ENV") != "docker" {
		constant.VideoAddr = getLocalAddr(constant.VideoAddr)
		constant.UserAddr = getLocalAddr(constant.UserAddr)
		constant.RelationAddr = getLocalAddr(constant.RelationAddr)
		constant.MessageAddr = getLocalAddr(constant.MessageAddr)
		constant.FavoriteAddr = getLocalAddr(constant.FavoriteAddr)
		constant.CommentAddr = getLocalAddr(constant.CommentAddr)
	}
	log.Println("viper读取配置文件成功")
}

func getLocalAddr(addr string) string {
	e := strings.Split(addr, ":")
	return "127.0.0.1:" + e[1]
}
