package model

type Favorite struct {
	ID      uint64 `gorm:"primaryKey" json:"id"`
	UserID  uint64 `gorm:"not null;uniqueIndex:idx_user_video" json:"user_id"`
	VideoID uint64 `gorm:"not null;uniqueIndex:idx_user_video" json:"video_id"`
	//IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
	//CreatedTime time.Time `gorm:"not null" json:"created_time"`
}
