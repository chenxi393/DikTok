package model

import "time"

type Message struct {
	// 消息id
	ID int64 `json:"id"`
	// 消息内容
	Content string `json:"content"`
	// 消息发送时间 yyyy-MM-dd HH:MM:ss
	CreateTime time.Time `json:"create_time"`
	// 消息发送者id
	FromUserID uint64 `json:"from_user_id"`
	// 消息接收者id
	ToUserID uint64 `json:"to_user_id"`
}
