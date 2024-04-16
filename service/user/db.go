package main

import (
	"douyin/storage/database"
	"douyin/storage/database/model"

	"gorm.io/gorm"
)

func CreateUser(user *model.User) (int64, error) {
	err := database.DB.Model(&model.User{}).Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// 这里用事务 更新缓存 我们认为这个用户修改的行为 不大导致缓存不一致的情况
func UpdateUser(user *model.User) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		err := database.DB.Model(user).UpdateColumns(user).Error
		if err != nil {
			return err
		}
		return SetUserInfo(user)
	})
}

func SelectUserByName(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Model(&model.User{}).Where("username = ? ", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SelectUserByID(userID int64) (*model.User, error) {
	var user model.User
	err := database.DB.Model(&model.User{}).Where("id = ? ", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SelectWorkCount(userID int64) (int64, error) {
	var cnt int64
	err := database.DB.Model(&model.User{}).Select("work_count").Where("id = ? ", userID).First(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 通过一组id 批量获取用户信息
func SelectUserListByIDs(userIDs []int64) ([]model.User, error) {
	var users []model.User
	// (?)  ( ? )会多加一个括号
	err := database.DB.Model(&model.User{}).Where("id IN (?)  ", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
