package util

import (
	"douyin/config"

	"go.uber.org/zap"
)

func InitZap() {
	// 先打印到控制台吧
	var logger *zap.Logger
	if config.System.Mode == "debug" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer logger.Sync()
	zap.ReplaceGlobals(logger) //返回值似乎是一个取消函数
	logger.Info("zap初始化: 成功")
}
