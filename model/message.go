package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	Content    string    `gorm:"not null;type:text" json:"content"`
	CreateTime time.Time `gorm:"not null;index" json:"create_time"` // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID uint64    `gorm:"not null;index:idx_user_touser" json:"from_user_id"`
	ToUserID   uint64    `gorm:"not null;index:idx_user_touser" json:"to_user_id"`
}

type MessageMongo struct {
	Id_        primitive.ObjectID `bson:"_id,omitempty"`
	ID         int64              `gorm:"not null;type:id" json:"id" bson:"id"`
	Content    string             `gorm:"not null;type:text" json:"content" bson:"content"`
	CreateTime time.Time          `gorm:"not null;index" json:"create_time" bson:"create_time"` // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID uint64             `gorm:"not null;index:idx_user_touser" json:"from_user_id" bson:"from_user_id"`
	ToUserID   uint64             `gorm:"not null;index:idx_user_touser" json:"to_user_id" bson:"to_user_id"`
}
