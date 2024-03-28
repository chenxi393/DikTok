package mq

import (
	"douyin/config"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var connRabbitMQ *amqp.Connection
var ProduceChannel *amqp.Channel

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
	ProduceChannel, err = connRabbitMQ.Channel()
	if err != nil {
		zap.L().Error("创建channel失败", zap.Error(err))
	}
	// 下面这段代码 确保mq异常时可以重新连接
	go func() {
		defer ProduceChannel.Close()
		for {
			reason, ok := <-ProduceChannel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || ProduceChannel.IsClosed() {
				zap.L().Error("channel closed")
				_ = ProduceChannel.Close() // close again, ensure closed flag set when connection closed
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
					ProduceChannel = ch
					break
				}
				zap.L().Sugar().Errorf("channel recreate failed, err: %v", err)
			}
		}

	}()
}

func InitComment() {
	initMQ()
	err := ProduceChannel.ExchangeDeclare(
		"comment", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		zap.L().Sugar().Errorf("初始化comment exchange失败: %v", err)
		return
	}
}

func InitRelation() {
	initMQ()
	err := ProduceChannel.ExchangeDeclare(
		"relation", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		zap.L().Sugar().Errorf("初始化Relation exchange失败: %v", err)
		return
	}
}

func InitFavorite() {
	initMQ()
	err := ProduceChannel.ExchangeDeclare(
		"favorite", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		zap.L().Sugar().Errorf("初始化favorite exchange失败: %v", err)
		return
	}
}
