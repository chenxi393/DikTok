package database

import (
	"douyin/model"
)

func CreateUser(user *model.User) (uint64, error) {
	//TODO 待验证 不select指定更新的自动 无法使用默认值？？？？
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

func SelectWorkCount(userID uint64) (int64, error) {
	var cnt int64
	err := global_db.Model(&model.User{}).Select("work_count").Where("id = ? ", userID).First(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 通过一组id 批量获取用户信息
func SelectUserListByIDs(userIDs []uint64) ([]model.User, error) {
	var users []model.User
	// (?)  ( ? )会多加一个括号
	err := global_db.Model(&model.User{}).Where("id IN (?)  ", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
