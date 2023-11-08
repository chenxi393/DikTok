package util_test

import (
	"douyin/config"
	"douyin/package/util"
	"testing"

	"github.com/spf13/viper"
)

func TestChatGPT(t *testing.T) {
	viper.AddConfigPath("../../config/")
	config.Init()
	content := "介绍一下美国"
	util.RequestToSparkAPI(content)
}
