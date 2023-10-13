package database

import "douyin/model"

func FavoriteVideo(userID, videoID uint64) error {
	favorite := model.UserFavoriteVideo{
		UserID:  userID,
		VideoID: videoID,
	}
	return global_db.Model(&model.UserFavoriteVideo{}).Create(&favorite).Error
}

func UnFavoriteVideo(userID, videoID uint64) error {
	favorite := model.UserFavoriteVideo{
		UserID:  userID,
		VideoID: videoID,
	}
	return global_db.Model(&model.UserFavoriteVideo{}).Delete(&favorite).Error
}

func SelectFavoriteVideoByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := global_db.Model(&model.UserFavoriteVideo{}).Select("video_id").Where("user_id = ?", userID).Find(&res).Error
	return res, err
}

