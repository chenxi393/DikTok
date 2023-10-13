package database

import (
	"douyin/config"
	"log"
	"testing"

	"github.com/spf13/viper"
)

func TestXxx(t *testing.T) {
	viper.AddConfigPath("../config/")
	config.Init()
	InitMysql()
	var i int64
	res, _ := SelectFeedVideoList(30, &i)
	log.Printf("%#v\n", res)
}
