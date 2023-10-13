package model

type Follow struct{
	ID          uint64    `json:"id"`
	//IsDeleted   uint8     `gorm:"default:0;not null" json:"is_deleted"`
	UserID      uint64    `gorm:"not null" json:"user_id"`
	ToUserID    uint64    `gorm:"not null" json:"to_user_id"`
	//CreatedTime time.Time `gorm:"not null" json:"created_time"`
}
