package model

import "time"

type Comment struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	VideoID     uint64    `gorm:"not null;index" json:"video_id"`
	UserID      uint64    `gorm:"not null;index" json:"user_id"`
	Content     string    `gorm:"not null;type:varchar(255)" json:"content"`
	CreatedTime time.Time `gorm:"not null" json:"created_time"`
	//IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
}
