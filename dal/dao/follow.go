package dao

import "douyin/dal/model"

func SelectFollowingByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := global_db.Model(&model.Follow{}).Select("to_user_id").Where("user_id = ?", userID).Find(&res).Error
	return res, err
}
