package main

import (
	"context"
	"douyin/model"
	"douyin/package/util"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MessageClinet *mongo.Collection

func InitMongoDB() func() {
	// log monitor
	var logMonitor = event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			log.Printf("mongo reqId:%d start on db:%s cmd:%s sql:%+v", startedEvent.RequestID, startedEvent.DatabaseName,
				startedEvent.CommandName, startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			log.Printf("mongo reqId:%d exec cmd:%s success duration %d ns", succeededEvent.RequestID,
				succeededEvent.CommandName, succeededEvent.Duration)
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			log.Printf("mongo reqId:%d exec cmd:%s failed duration %d ns", failedEvent.RequestID,
				failedEvent.CommandName, failedEvent.Duration)
		},
	}

	//1.建立连接
	// TODO 待配置 暂时写死
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
		"mongodb://localhost:27017",
	).SetAuth(options.Credential{
		Username: "root",
		Password: "123456",
	}).SetMonitor(&logMonitor).SetConnectTimeout(5*time.Second))
	if err != nil {
		log.Panic(err)
	}

	//2.选择数据库
	db := client.Database("db")

	//3.选择表
	MessageClinet = db.Collection("message")

	return func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func CreateMessage(userID, toUSerId uint64, content string) error {
	uid, _ := util.GetSonyFlakeID()
	msg := model.MessageMongo{
		Content:    content,
		FromUserID: userID,
		ToUserID:   toUSerId,
		CreateTime: time.Now(),
		ID:         int64(uid),
	}
	_, err := MessageClinet.InsertOne(context.Background(), msg)
	if err != nil {
		return err
	}
	return nil
}

func GetMessages(userID, toUserID uint64, msgTime int64) ([]model.MessageMongo, error) {
	newMsgTime := time.UnixMilli(msgTime)
	msgs := make([]model.MessageMongo, 0)

	// 构建排序规则
	sort := bson.D{{Key: "create_time", Value: -1}}
	// TODO 这里是不是可以用联合索引 create_time
	cursor, err := MessageClinet.Find(context.TODO(), bson.M{
		"from_user_id": bson.M{"$in": []uint64{userID, toUserID}},
		"to_user_id":   bson.M{"$in": []uint64{userID, toUserID}},
		"create_time":  bson.M{"$gt": newMsgTime},
	}, options.Find().SetSort(sort))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.Background(), &msgs)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// 用来呈现好友列表的第一条消息
func GetNewestMessage(userID, toUserID uint64) (*model.MessageMongo, error) {
	msg := model.MessageMongo{}

	sort := bson.D{{Key: "create_time", Value: -1}}
	err := MessageClinet.FindOne(context.TODO(), bson.M{
		"from_user_id": bson.M{"$in": []uint64{userID, toUserID}},
		"to_user_id":   bson.M{"$in": []uint64{userID, toUserID}},
	}, options.FindOne().SetSort(sort)).Decode(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
