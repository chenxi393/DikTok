package dao

import "douyin/dal/model"

func SelectFavoriteVideoByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := global_db.Model(&model.UserFavoriteVideo{}).Select("video_id").Where("user_id = ?", userID).Find(&res).Error
	return res, err
}
