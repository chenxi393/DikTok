package dao

import (
	"douyin/dal/model"

	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao() *UserDao {
	return &UserDao{global_db}
}

func (dao *UserDao) SelectUserByName(username string) (*model.User, error) {
	var user model.User
	err := dao.Model(&model.User{}).Where("username = ? ", username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) CreateUser(user *model.User) (uint64, error) {
	// 不select指定更新的自动 无法使用默认值？？？？ TODO 待验证
	err := dao.Model(&model.User{}).Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
