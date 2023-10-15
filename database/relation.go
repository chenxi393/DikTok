package database

import (
	"douyin/model"

	"gorm.io/gorm"
)

func Follow(userID, toUserID uint64, cnt int64) error {
	follow := model.Follow{
		UserID:   userID,
		ToUserID: toUserID,
	}
	return global_db.Transaction(func(tx *gorm.DB) error {
		// TODO前面很多POST应该都需要检查一下
		// 关注表里更新
		var err error
		var ff model.Follow
		if cnt == 1 {
			// FIX没有唯一约束 可以插入多条记录 要不插入前检查一下
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
		return nil
	})
}

func SelectFollowingByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64,0)
	err := global_db.Model(&model.Follow{}).Select("to_user_id").Where("user_id = ?", userID).Find(&res).Error
	return res, err
}

func SelectFollowerByUserID(userID uint64) ([]uint64, error) {
	res := make([]uint64, 0)
	err := global_db.Model(&model.Follow{}).Select("user_id").Where("to_user_id = ?", userID).Find(&res).Error
	return res, err
}
