package database

import (
	"douyin/model"
)

// CreateVideo 新增视频，返回的 videoID 是为了将 videoID 放入布隆过滤器
// FIX 这里需要事务（一个不成功要回滚） 要更新用户的视频数+1 也更新到redis
// 这里简单的先写到数据库 后序使用redis+ 布隆过滤器
func CreateVideo(video *model.Video) (uint64, error) {
	// 这里要用事务代替
	err := global_db.Model(&model.Video{}).Save(video).Error
	if err != nil {
		return 0, err
	}
	cnt, err := SelectWorkCount(video.AuthorID)
	if err != nil {
		//FIX 这里要回退
		return 0, err
	}
	err = global_db.Model(&model.User{}).Where("id = ?", video.AuthorID).Update("work_count", cnt+1).Error
	if err != nil {
		//FIX 这里要回退
		return video.ID, err
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
