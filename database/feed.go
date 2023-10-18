package database

import (
	"douyin/model"
	"douyin/response"
	"time"
)

// Scan支持的数据类型仅为struct及struct slice以及它们的指针类型
// Scan要不结构体加tag  gorm:"column:col_name" 指定列名 要不改造结构体
func SelectFeedVideoList(numberVideos int, lastTime *int64) ([]response.VideoData, error) {
	if lastTime == nil || *lastTime == 0 { //很多时候这种写法会有空指针的问题
		// 所以得产生一个新的变量 杜绝空指针问题 内存回收交给gc
		currentTime := time.Now().UnixMilli()
		lastTime = &currentTime
	}
	// TODO 下面的时间要用小于 可以考虑减1 用小于等于（为了使用索引？？）
	res := make([]response.VideoData, 0, 30)
	// 这里使用外连接 双表联查 可以考虑改多次单表 联查太麻烦
	err := global_db.Model(&model.User{}).Select(`user.*,
    video.id as vid,
    video.play_url,
    video.cover_url,
    video.favorite_count as vfavorite_count,
    video.comment_count,
    video.title`).Joins(
		"right join video on video.author_id = user.id").Where("video.publish_time < ? ",
		time.UnixMilli(*lastTime)).Order("video.publish_time desc").Limit(numberVideos).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectPublishTimeByVideoID(vid uint64) (time.Time, error) {
	var publishTime time.Time
	err := global_db.Model(&model.Video{}).Select("publish_time").Where("id = ?", vid).First(&publishTime).Error
	return publishTime, err
}
