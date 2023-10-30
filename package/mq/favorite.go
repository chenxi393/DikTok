package mq

import (
	"context"
	"douyin/database"
	"fmt"
	"strconv"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func initFavorite() {
	err := produceChannel.ExchangeDeclare(
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
	go favoriteConsume()
}

// flag:1 表示点赞 -1表示取消赞
func SendFavoriteMessage(userID, videoID uint64, flag int) error {
	err := produceChannel.PublishWithContext(
		context.Background(),
		"favorite",
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //开启消息持久化
			ContentType:  "text/plain",
			Body:         []byte(fmt.Sprintf("%d:%d:%d", userID, videoID, flag)),
		},
	)
	if err != nil {
		zap.L().Error("发布消息失败", zap.Error(err))
	}
	return err
}

func favoriteConsume() {
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
		"favorite", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		zap.L().Sugar().Errorf("favorite 创建exchange失败 : %v", err)
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
		q.Name,     // queue name
		"",         // routing key
		"favorite", // exchange
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
		// 注意要是多次消费不成功 仍入死信队列？？ 配置最多消费次数
		// TODO 应该可以配置一条消息最多的消费次数
		// 以及消费失败的罚时
		for message := range msgs {
			// For example, just show received message in console.
			zap.L().Sugar().Infof("Received message: %s\n", message.Body)
			ans := strings.Split(string(message.Body), ":")
			userID, err := strconv.ParseUint(ans[0], 10, 64)
			if err != nil {
				zap.L().Sugar().Error(err)
				message.Ack(true)
				continue
			}
			videoID, err := strconv.ParseUint(ans[1], 10, 64)
			if err != nil {
				zap.L().Sugar().Error(err)
				message.Ack(true)
				continue
			}
			cnt, _ := strconv.ParseInt(ans[2], 10, 64)
			err = database.FavoriteVideo(userID, videoID, cnt)
			if err != nil {
				zap.L().Sugar().Error(err)
				message.Ack(true)
				continue
			}
			message.Ack(false) //手动确认消息
		}
	}()
	// Close the channel.
	<-forever
}
