package util

import (
	"douyin/config"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// 获取七牛云的上传凭证 有效时间为10 min
func GetUploadToken(fileName string) string {
	putPolicy := storage.PutPolicy{
		Scope:           config.System.Qiniu.Bucket + ":" + fileName,
		IsPrefixalScope: 1,
		Expires:         uint64(time.Now().Unix()) + 600, // 给了10min 给用户上传
	}
	mac := qbox.NewMac(config.System.Qiniu.AccessKey, config.System.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}
