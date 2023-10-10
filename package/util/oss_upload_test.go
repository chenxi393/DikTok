package util_test

import (
	"context"
	"douyin/config"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestUploadVideoToOSS(t *testing.T) {
	bucket := "chenxi-douyin"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	viper.AddConfigPath("../../config/")
	config.Init()
	
	log.Println(config.SystemConfig.AccessKey, config.SystemConfig.SecretKey)
	mac := qbox.NewMac(config.SystemConfig.AccessKey, config.SystemConfig.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	text := "Hello, world!fsuhdiuahdojuahiusd"
	reader := strings.NewReader(text)
	reader.Len()
	datalen := int64(len(text))
	err := formUploader.Put(context.Background(), &ret, upToken, "测试使用.txt", reader, datalen, &putExtra)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}
