package mq

import (
	"douyin/config"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var connRabbitMQ *amqp.Connection
var produceChannel *amqp.Channel

func initMQ() {
	var err error
	// Create a new RabbitMQ connection.
	connRabbitMQ, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.System.MQ.User,
		config.System.MQ.Password,
		config.System.MQ.Host,
		config.System.MQ.Port))
	if err != nil {
		zap.L().Fatal("RabbitMQ连接: 失败", zap.Error(err))
	}
	zap.L().Info("RabbitMQ连接: 成功")
	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	produceChannel, err = connRabbitMQ.Channel()
	if err != nil {
		zap.L().Error("创建channel失败", zap.Error(err))
	}
	// 不能使用defer关闭
	// defer produceChannel.Close()

	// 下面这段代码 确保mq异常时可以重新连接
	go func() {
		for {
			reason, ok := <-produceChannel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || produceChannel.IsClosed() {
				zap.L().Error("channel closed")
				_ = produceChannel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			zap.L().Sugar().Warnf("channel closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(3 * time.Second)

				ch, err := connRabbitMQ.Channel()
				if err == nil {
					zap.L().Sugar().Info("channel recreate success")
					produceChannel = ch
					break
				}
				zap.L().Sugar().Errorf("channel recreate failed, err: %v", err)
			}
		}
	}()
}
