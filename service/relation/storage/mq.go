package storage

import (
	"context"
	"fmt"

	"diktok/storage/mq"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// t表示类型 1关注 -1取关
func SendFollowMessage(userID, toUserID int64, t int) {
	// Attempt to publish a message to the queue.
	err := mq.ProduceChannel.PublishWithContext(
		context.Background(),
		"relation", // 默认的exchange
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //开启消息持久化
			ContentType:  "text/plain",
			Body:         []byte(fmt.Sprintf("%d:%d:%d", userID, toUserID, t)),
		},
	)
	if err != nil {
		zap.L().Error("发布消息失败", zap.Error(err))
	}
}
