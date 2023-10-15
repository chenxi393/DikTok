package database

import (
	"douyin/model"

	"gorm.io/gorm"
)

// CreateVideo 新增视频，返回的 videoID 是为了将 videoID 放入布隆过滤器
// 这里简单的先写到数据库 后序使用redis + 布隆过滤器
func CreateVideo(video *model.Video) (uint64, error) {
	err := global_db.Transaction(func(tx *gorm.DB) error {
		// If value doesn't contain a matching primary key, value is inserted.
		err := tx.Model(&model.Video{}).Save(video).Error
		if err != nil {
			return err
		}
		cnt, err := SelectWorkCount(video.AuthorID)
		if err != nil {
			return err
		}
		err = tx.Model(&model.User{}).Where("id = ?", video.AuthorID).Update("work_count", cnt+1).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return video.ID, nil
}

func SelectVideosByUserID(userID uint64) ([]model.Video, error) {
	videos := make([]model.Video, 0)
	err := global_db.Model(&model.Video{}).Where("author_id = ? ", userID).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 根据视频ID集合查询视频信息 批量查询
func SelectVideoListByVideoID(videoIDList []uint64) ([]model.Video, error) {
	res := make([]model.Video, 0, len(videoIDList))
	// 这里按照id倒叙 其实id就能保证时间顺序了
	err := global_db.Where("id IN (?)", videoIDList).Order("id desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, err
}
