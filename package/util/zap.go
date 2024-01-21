package util

import (
	"context"
	"douyin/config"
	"douyin/package/constant"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func InitZap() {
	// 先打印到控制台吧
	var logger *zap.Logger
	if config.System.Mode == constant.DebugMode {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	defer logger.Sync()
	zap.ReplaceGlobals(logger) //返回值似乎是一个取消函数

	// TODO 还没改造 otel日志 将log 嵌入 span里面 
	otelLogger := otelzap.New(zap.NewExample())
	defer otelLogger.Sync()

	undo := otelzap.ReplaceGlobals(otelLogger)
	defer undo()
	otelzap.L().Info("replaced zap's global loggers")
	otelzap.Ctx(context.TODO()).Info("... and with context")
	
	logger.Info("zap初始化: 成功")
}
