package main

// import (
// 	"context"
// 	"diktok/package/mq"
// 	"fmt"

// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"go.uber.org/zap"
// )

// // flag:1 表示点赞 -1表示取消赞
// func SendFavoriteMessage(userID, videoID uint64, flag int) error {
// 	err := mq.ProduceChannel.PublishWithContext(
// 		context.Background(),
// 		"favorite",
// 		"",
// 		false,
// 		false,
// 		amqp.Publishing{
// 			DeliveryMode: amqp.Persistent, //开启消息持久化
// 			ContentType:  "text/plain",
// 			Body:         []byte(fmt.Sprintf("%d:%d:%d", userID, videoID, flag)),
// 		},
// 	)
// 	if err != nil {
// 		zap.L().Error("发布消息失败", zap.Error(err))
// 	}
// 	return err
// }
