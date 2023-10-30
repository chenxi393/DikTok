package model

type Follow struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	UserID   uint64 `gorm:"not null;uniqueIndex:idx_user_touser" json:"user_id"`
	ToUserID uint64 `gorm:"not null;uniqueIndex:idx_user_touser;index" json:"to_user_id"`
	//IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
	//CreatedTime time.Time `gorm:"not null" json:"created_time"`
}
