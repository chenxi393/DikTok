package database

import (
	"douyin/model"
	"douyin/package/constant"
	"douyin/storage/cache"

	"gorm.io/gorm"
)

func Follow(userID, toUserID uint64, cnt int64) error {
	follow := model.Follow{
		UserID:   userID,
		ToUserID: toUserID,
	}
	return constant.DB.Transaction(func(tx *gorm.DB) error {
		// 关注表里更新
		var err error
		var ff model.Follow
		if cnt == 1 {
			//  这里设置联合唯一索引 应该不需要检查了
			err = tx.Model(&model.Follow{}).Create(&follow).Error
		} else if cnt == -1 {
			err = tx.Model(&model.Follow{}).Where("user_id = ? AND to_user_id = ?", userID, toUserID).Delete(&ff).Error
		}
		if err != nil {
			return err
		}
		// 然后更新用户的关注数
		user := model.User{ID: userID}
		// Model会检查主键
		err = tx.Model(&user).First(&user).Error
		if err != nil {
			return err
		}
		err = tx.Model(&user).Update("follow_count", user.FollowCount+cnt).Error
		if err != nil {
			return err
		}
		// 被关注用户的被关注数+1
		toUser := model.User{ID: toUserID}
		err = tx.Model(&toUser).First(&toUser).Error
		if err != nil {
			return err
		}
		err = tx.Model(&toUser).Update("follower_count", toUser.FollowerCount+cnt).Error
		if err != nil {
			return err
		}
		return cache.FollowAction(userID, toUserID, cnt)
	})
}

func SelectFollowingByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := constant.DB.Model(&model.Follow{}).Select("to_user_id").Where("user_id = ?", userID).Order("id desc").Find(&res).Error
	return res, err
}

func SelectFollowerByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := constant.DB.Model(&model.Follow{}).Select("user_id").Where("to_user_id = ?", userID).Order("id desc").Find(&res).Error
	return res, err
}

// 查询userID 有没有关注 id
func IsFollowed(userID uint64, id uint64) (bool, error) {
	var cnt int64
	err := constant.DB.Model(&model.Follow{}).Where("user_id= ? AND to_user_id = ? ", userID, id).Count(&cnt).Error
	if err != nil {
		return false, err
	} else if cnt == 0 {
		return false, nil
	}
	return true, nil
}
