// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameFollow = "follow"

// Follow mapped from table <follow>
type Follow struct {
	ID       int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID   int64 `gorm:"column:user_id;not null" json:"user_id"`
	ToUserID int64 `gorm:"column:to_user_id;not null" json:"to_user_id"`
}

// TableName Follow's table name
func (*Follow) TableName() string {
	return TableNameFollow
}