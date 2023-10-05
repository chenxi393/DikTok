package config_test

import (
	"douyin/config"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	viper.AddConfigPath(".")
	config.Init()
	fmt.Fprintf(os.Stderr, "%#v", config.SystemConfig)
}
