// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameCommentContent = "comment_content"

// CommentContent mapped from table <comment_content>
type CommentContent struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Content   string         `gorm:"column:content;not null;comment:评论内容 后续可以考虑垂直分出去" json:"content"` // 评论内容 后续可以考虑垂直分出去
	Extra     string         `gorm:"column:extra;not null;comment:回复用户 @用户[] 评论图片等" json:"extra"`     // 回复用户 @用户[] 评论图片等
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;not null" json:"deleted_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName CommentContent's table name
func (*CommentContent) TableName() string {
	return TableNameCommentContent
}
