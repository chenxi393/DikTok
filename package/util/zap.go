package util

import (
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
	// 如何做到 log 和trace 的联动呢？？ 不靠附在span里？？
	// zap.L()   otelzap.L().Ctx(c.UserContext()) 做一个全局的替换
	otelLogger := otelzap.New(logger)
	defer otelLogger.Sync()

	otelzap.ReplaceGlobals(otelLogger)
	//defer undo()

	logger.Info("zap初始化: 成功")
}
