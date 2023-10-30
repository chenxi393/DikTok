package model

import "time"

type Message struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	Content    string    `gorm:"not null;" json:"content"`
	CreateTime time.Time `gorm:"not null;index" json:"create_time"` // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID uint64    `gorm:"not null;index:idx_user_touser" json:"from_user_id"`
	ToUserID   uint64    `gorm:"not null;index:idx_user_touser" json:"to_user_id"`
}
