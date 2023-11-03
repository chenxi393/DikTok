package model

import "time"

// 这两个count是不是可以考虑解耦 要的时候再去计算
// 23.11.03 Title 增加全文索引 以便于搜索 ngram全文索引支持中文的插件
type Video struct {
	ID            uint64    `json:"id"`
	AuthorID      uint64    `gorm:"not null;index" json:"author_id"`
	PlayURL       string    `gorm:"type:varchar(255);not null" json:"play_url"`
	CoverURL      string    `gorm:"type:varchar(255);not null" json:"cover_url"`
	Title         string    `gorm:"type:varchar(63);index:,class:FULLTEXT,option:WITH PARSER ngram;not null" json:"title"`
	PublishTime   time.Time `gorm:"not null;index" json:"publish_time"`
	FavoriteCount int64     `gorm:"default:0;not null" json:"favorite_count"`
	CommentCount  int64     `gorm:"default:0;not null" json:"comment_count"`
	// 视频的分类 23.11.03新增 TODO 
	Topic string `gorm:"type:varchar(15);index:;not null" json:"topic"`
}
