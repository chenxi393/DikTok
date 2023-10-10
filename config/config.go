package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type MysqlConfig struct {
	// 需要解析的字段必须大写
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	UserName    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	MaxOpenConn int    `mapstructure:"maxOpenConn"`
	MaxIdleConn int    `mapstructure:"maxIdleConn"`
}

type HttpConfig struct {
	Host          string `mapstructure:"host"`
	Port          string `mapstructure:"port"`
	VideoAddress  string `mapstructure:"videoAddress"`
	ImageAddress  string `mapstructure:"imageAddress"`
	DefaltImagURL string `mapstructure:"defaltCoverURL"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
	Password string `mapstructure:"password"`
}

type System struct {
	MysqlMaster MysqlConfig `mapstructure:"mysqlMaster"`
	MysqlSlave  MysqlConfig `mapstructure:"mysqlSlave"`
	Mode        string      `mapstructure:"mode"`
	HttpAddress HttpConfig  `mapstructure:"httpAddress"`
	Redis       RedisConfig `mapstructure:"userRedis"`
	JwtSecret   string      `mapstructure:"jwtSecret"`
	UploadModel string      `mapstructure:"uploadModel"`
	AccessKey   string      `mapstructure:"accessKey"`
	SecretKey   string      `mapstructure:"secretKey"`
	Bucket      string      `mapstructure:"bucket"`
	OssDomain   string      `mapstructure:"ossDomain"`
}

var SystemConfig System

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal("fatal error config file: ", err.Error())
	}

	// 解析mysql的配置文件
	err = viper.Unmarshal(&SystemConfig)
	if err != nil {
		log.Fatal("fatal error unmarshal config: ", err.Error())
	}

	// 监视配置文件的变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件被修改")
	})
}
