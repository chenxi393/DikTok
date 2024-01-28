package main

import (
	"douyin/config"
	"testing"

	"github.com/spf13/viper"
)

func TestChatGPT(t *testing.T) {
	viper.AddConfigPath("../../config/")
	config.Init()
	content := "介绍一下美国"
	requestToSparkAPI(content)
}
