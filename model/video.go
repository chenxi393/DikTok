package model

import "time"
//pubulish字段是不是建一个索引比较好
type Video struct {
	ID            uint64    `json:"id"`
	PublishTime   time.Time `gorm:"not null" json:"publish_time"`
	AuthorID      uint64    `gorm:"not null" json:"author_id"`
	PlayURL       string    `gorm:"type:varchar(777);not null" json:"play_url"`
	CoverURL      string    `gorm:"type:varchar(777);not null" json:"cover_url"`
	FavoriteCount int64     `gorm:"default:0;not null" json:"favorite_count"`
	CommentCount  int64     `gorm:"default:0;not null" json:"comment_count"`
	Title         string    `gorm:"type:varchar(63);not null" json:"title"`
}