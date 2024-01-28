package database

import (
	"douyin/model"
	"douyin/package/constant"
)

func CreateUser(user *model.User) (uint64, error) {
	err := constant.DB.Model(&model.User{}).Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func UpdateUser(user *model.User) error {
	err := constant.DB.Model(user).UpdateColumns(user).Error
	return err
}

func SelectUserByName(username string) (*model.User, error) {
	var user model.User
	err := constant.DB.Model(&model.User{}).Where("username = ? ", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SelectUserByID(userID uint64) (*model.User, error) {
	var user model.User
	err := constant.DB.Model(&model.User{}).Where("id = ? ", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SelectWorkCount(userID uint64) (int64, error) {
	var cnt int64
	err := constant.DB.Model(&model.User{}).Select("work_count").Where("id = ? ", userID).First(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 通过一组id 批量获取用户信息
func SelectUserListByIDs(userIDs []uint64) ([]model.User, error) {
	var users []model.User
	// (?)  ( ? )会多加一个括号
	err := constant.DB.Model(&model.User{}).Where("id IN (?)  ", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
