package mq

import (
	"context"
	"douyin/model"
	"douyin/storage/database"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func InitComment() {
	initMQ()
	err := produceChannel.ExchangeDeclare(
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
	go commentConsume()
}

func SendCommentMessage(msg *model.Comment) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error(err.Error())
	}
	err = produceChannel.PublishWithContext(
		context.Background(),
		"comment",
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //开启消息持久化
			ContentType:  "text/plain",
			Body:         msgJSON,
		},
	)
	if err != nil {
		zap.L().Error("发布消息失败", zap.Error(err))
	}
	return err
}

func commentConsume() {
	// Open a new channel.
	channel, err := connRabbitMQ.Channel()
	if err != nil {
		zap.L().Fatal("打开rabbitMQ channel失败", zap.Error(err))
	}
	defer channel.Close()
	// 设置消息队列分发消息的测量 一个线程不多于2个
	_ = channel.Qos(
		2,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	err = channel.ExchangeDeclare(
		"comment", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		zap.L().Sugar().Errorf("comment 创建exchange失败 : %v", err)
	}
	// 创建一个无名队列 临时的
	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		zap.L().Sugar().Error(err, "Failed to declare a queue")
	}
	err = channel.QueueBind(
		q.Name,    // queue name
		"",        // routing key
		"comment", // exchange
		false,
		nil,
	)
	if err != nil {
		zap.L().Sugar().Error(err, "Failed to bind a queue")
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		zap.L().Error("Failed to register a consumer", zap.Error(err))
	}
	// Open a channel to receive messages.
	forever := make(chan struct{})
	go func() {
		for message := range msgs {
			// For example, just show received message in console.
			zap.L().Sugar().Infof("Received message")
			comment := model.Comment{}
			err := json.Unmarshal(message.Body, &comment)
			if err != nil {
				zap.L().Error(err.Error())
				break
			}
			database.CommentAdd(&comment)
			message.Ack(false) //手动确认消息
		}
	}()
	// Close the channel.
	<-forever
}
