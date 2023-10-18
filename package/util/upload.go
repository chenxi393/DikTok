package util

import (
	"bytes"
	"context"
	"douyin/config"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
)

func UploadVideo(file []byte, fileName string) (string, string, error) {
	err := os.MkdirAll(config.SystemConfig.HttpAddress.VideoAddress, os.ModePerm)
	if err != nil {
		zap.L().Error(err.Error())
		return "", "", err
	}
	// 还得有个变量是宿主机ip
	path := "http://" + config.SystemConfig.MyIP + ":" + config.SystemConfig.HttpAddress.Port
	outputFilePath := filepath.Join(config.SystemConfig.HttpAddress.VideoAddress, fileName)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		zap.L().Error(err.Error())
		return "", "", err
	}
	defer outputFile.Close()
	_, err = outputFile.Write(file)
	if err != nil {
		zap.L().Error(err.Error())
		return "", "", err
	}
	zap.L().Info(fileName + "已成功写入文件夹")
	videoURL := path + "/video/" + fileName
	// if config.SystemConfig.Mode == "debug" { // 本地没有装ffmpeg 这里直接返回默认的url
	// 	return videoURL, config.SystemConfig.HttpAddress.DefaultCoverURL, nil
	// }
	err = GetVideoFrame(outputFilePath, fileName)
	if err != nil {
		return videoURL, config.SystemConfig.HttpAddress.DefaultCoverURL, nil
	}
	return videoURL, path + "/image/" + fileName + ".jpeg", nil
}

func GetVideoFrame(outputFilePath, fileName string) error {
	// 提取视频的第一帧
	picBuffer := bytes.NewBuffer(nil)
	err := ffmpeg.Input(outputFilePath).
		Filter("select", ffmpeg.Args{"gte(n,1)"}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(picBuffer).
		Run()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	img, err := imaging.Decode(picBuffer)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	// 创建图片文件夹
	err = os.MkdirAll(config.SystemConfig.HttpAddress.ImageAddress, os.ModePerm)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	// 保存
	imagPath := config.SystemConfig.HttpAddress.ImageAddress +"/"+ fileName + ".jpeg"
	err = imaging.Save(img, imagPath)
	if err != nil {
		zap.L().Error(err.Error())
		return nil
	}
	return nil
}

func UploadToOSS(fileName, filePath string) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: config.SystemConfig.Bucket,
	}
	mac := qbox.NewMac(config.SystemConfig.AccessKey, config.SystemConfig.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
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
	// 这里其实耗时很久 感觉有3/4秒
	err := formUploader.PutFile(context.Background(), &ret, upToken, fileName, filePath, &putExtra)
	if err != nil {
		zap.L().Error(err.Error())
		return "", nil
	}
	// FIX 10.11这里很奇怪的是 生成的URL 一天了都可以访问？？ 如果那就不用定时更新数据库的URL了
	// 如果采用私有空间 url有限制时间访问的 那就应该异步定时更新数据库的URL 否则用户访问不到啊
	// 这里好像是返回了CND 因为那个域名是开启了CDN的 都能访问
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	return storage.MakePrivateURL(mac, config.SystemConfig.OssDomain, ret.Key, deadline), nil
}
