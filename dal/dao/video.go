package dao

import "douyin/dal/model"

func SelectFavoriteVideoByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := global_db.Model(&model.UserFavoriteVideo{}).Select("video_id").Where("user_id = ?", userID).Find(&res).Error
	return res, err
}

// CreateVideo 新增视频，返回的 videoID 是为了将 videoID 放入布隆过滤器
// FIX 这里需要事务（一个不成功要回滚） 要更新用户的视频数+1 也更新到redis
// 这里简单的先写到数据库 后序使用redis+ 布隆过滤器
func CreateVideo(video *model.Video) (uint64, error) {
	err := global_db.Model(&model.Video{}).Save(video).Error
	if err != nil {
		return 0, err
	}
	cnt, err := SelectWorkCount(video.AuthorID)
	if err != nil {
		//FIX 这里要回退
		return 0, err
	}
	err = global_db.Model(&model.User{}).Where("id = ?",video.AuthorID).Update("work_count", cnt+1).Error
	if err != nil {
		//FIX 这里要回退
		return video.ID, err
	}
	return video.ID, nil
}
