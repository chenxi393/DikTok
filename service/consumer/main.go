package main

// import (
// 	"diktok/config"
// 	"diktok/storage/database"
// 	"diktok/package/util"
// 	"fmt"

// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"go.uber.org/zap"
// )

// var connRabbitMQ *amqp.Connection

// func main() {
// 	config.Init()
// 	util.InitZap()
// 	database.InitMySQL()
// 	var err error
// 	// Create a new RabbitMQ connection.
// 	connRabbitMQ, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
// 		config.System.MQ.User,
// 		config.System.MQ.Password,
// 		config.System.MQ.Host,
// 		config.System.MQ.Port))
// 	if err != nil {
// 		zap.L().Fatal("RabbitMQ连接: 失败", zap.Error(err))
// 	}
// 	zap.L().Info("RabbitMQ连接: 成功")

// 	CommentConsume()
// }
