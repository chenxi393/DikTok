package model

import "time"

// 系统通知 当关注
type Notice struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"not null" json:"content"`
	CreateTime time.Time `gorm:"not null;index" json:"create_time"` // 消息发送时间 yyyy-MM-dd HH:MM:ss
	UserID     int64     `gorm:"not null;index:idx_user_touser" json:"user_id"`
	HasRead    int       `gorm:"not null" json:"has_read"`
}
