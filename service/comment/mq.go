package main

// import (
// 	"context"
// 	"diktok/storage/database/model"
// 	"diktok/package/mq"
// 	"encoding/json"

// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"go.uber.org/zap"
// )

// func SendCommentMessage(msg *model.Comment) error {
// 	msgJSON, err := json.Marshal(msg)
// 	if err != nil {
// 		zap.L().Error(err.Error())
// 	}
// 	err = mq.ProduceChannel.PublishWithContext(
// 		context.Background(),
// 		"comment",
// 		"",
// 		false,
// 		false,
// 		amqp.Publishing{
// 			DeliveryMode: amqp.Persistent, //开启消息持久化
// 			ContentType:  "text/plain",
// 			Body:         msgJSON,
// 		},
// 	)
// 	if err != nil {
// 		zap.L().Error("发布消息失败", zap.Error(err))
// 	}
// 	return err
// }
