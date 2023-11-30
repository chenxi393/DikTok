package model

// 这些count是不是可以考虑去除 要的时候再去表里计算TODO
// 这里的头像和背景存储文件名而不是URL 方便更换OSS域名
type User struct {
	ID              uint64 `gorm:"primaryKey" json:"id"`
	Username        string `gorm:"uniqueIndex;type:varchar(63);not null" json:"username"`
	Password        string `gorm:"type:varchar(255);not null" json:"password"`
	Avatar          string `gorm:"type:varchar(255);not null" json:"avatar"`
	BackgroundImage string `gorm:"type:varchar(255);not null" json:"background_image"`
	Signature       string `gorm:"type:varchar(255);not null" json:"signature"`
	FollowCount     int64  `gorm:"default:0;not null" json:"follow_count"`
	FollowerCount   int64  `gorm:"default:0;not null" json:"follower_count"`
	TotalFavorited  int64  `gorm:"default:0;not null" json:"total_favorited"`
	FavoriteCount   int64  `gorm:"default:0;not null" json:"favorite_count"`
	WorkCount       int64  `gorm:"default:0;not null" json:"work_count"`
}
