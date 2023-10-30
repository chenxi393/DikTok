package model

import "time"

// 这两个count是不是可以考虑解耦 要的时候再去计算
type Video struct {
	ID            uint64    `json:"id"`
	AuthorID      uint64    `gorm:"not null;index" json:"author_id"`
	PlayURL       string    `gorm:"type:varchar(777);not null" json:"play_url"`
	CoverURL      string    `gorm:"type:varchar(777);not null" json:"cover_url"`
	Title         string    `gorm:"type:varchar(63);not null" json:"title"`
	PublishTime   time.Time `gorm:"not null;index" json:"publish_time"`
	FavoriteCount int64     `gorm:"default:0;not null" json:"favorite_count"`
	CommentCount  int64     `gorm:"default:0;not null" json:"comment_count"`
}
