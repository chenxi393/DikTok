package model

import "time"

// 这两个count是不是可以考虑删除 要的时候再去计算
// 23.11.03 Title 增加全文索引 以便于搜索 ngram全文索引支持中文的插件 默认分词2
// 存URL存在更换CDN域名之类的问题 可以考虑存文件名 然后灵活更换
type Video struct {
	ID            uint64    `json:"id"`
	AuthorID      uint64    `gorm:"not null;index" json:"author_id"`
	PlayURL       string    `gorm:"type:varchar(255);not null" json:"play_url"`
	CoverURL      string    `gorm:"type:varchar(255);not null" json:"cover_url"`
	Title         string    `gorm:"type:varchar(63);index:idx_title_topic,class:FULLTEXT,option:WITH PARSER ngram;not null" json:"title"`
	PublishTime   time.Time `gorm:"not null;index" json:"publish_time"`
	FavoriteCount int64     `gorm:"default:0;not null" json:"favorite_count"`
	CommentCount  int64     `gorm:"default:0;not null" json:"comment_count"`
	// 视频的分类 23.11.03新增 前两个为固定字段 后面为tag隐式搜索
	Topic string `gorm:"type:varchar(63);index:idx_title_topic,class:FULLTEXT,option:WITH PARSER ngram;not null" json:"topic"`
}
