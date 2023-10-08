package dao

import (
	"douyin/dal/model"
)

func CreateUser(user *model.User) (uint64, error) {
	// 不select指定更新的自动 无法使用默认值？？？？ TODO 待验证
	err := global_db.Model(&model.User{}).Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func SelectUserByName(username string) (*model.User, error) {
	var user model.User
	err := global_db.Model(&model.User{}).Where("username = ? ", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SelectUserByID(userID uint64) (*model.User, error) {
	var user model.User
	err := global_db.Model(&model.User{}).Where("id = ? ", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 查询userID 有没有关注 id
func IsFollowed(userID uint64, id uint64) (bool, error) {
	var cnt int64
	err := global_db.Model(&model.Follow{}).Where("user_id= ? AND to_user_id = ? ", userID, id).Count(&cnt).Error
	if err != nil {
		return false, err
	} else if cnt == 0 {
		return false, nil
	}
	return true, nil
}
