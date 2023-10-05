package service

import (
	"douyin/dal/dao"
	"douyin/dal/model"
	"douyin/package/util"
	"douyin/response"
	"fmt"

	"go.uber.org/zap"
)

type UserService struct {
	// 密码，最长32个字符
	Password string `json:"password"`
	// 注册用户名，最长32个字符
	Username string `json:"username"`
}

func (service UserService) RegisterService() (*response.UserRegisterOrLogin, error) {
	logTag := "service.user.Register err:"
	// 这个复杂度判断是不是放在上一层比较好？？
	if len(service.Password) < 6 || len(service.Password) > 32 {
		return nil, fmt.Errorf("密码格式错误")
	}
	if service.Password == "123456" {
		return nil, fmt.Errorf("密码太简单")
	}

	userDao := dao.NewUserDao()
	//先判断用户存不存在
	user, err := userDao.SelectUserByName(service.Username)
	if err != nil {
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	if user.ID != 0 {
		err := fmt.Errorf("用户名已被注册")
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	// 对密码进行加密并存储
	encryptedPassword := util.BcryptHash(service.Password)
	user = &model.User{
		Username: service.Username,
		Password: encryptedPassword,
	}
	userID, err := userDao.CreateUser(user)
	if err != nil {
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	// 签发token
	token, err := util.SignToken(userID)
	if err != nil {
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	// 将 UserID 添加到布隆过滤器中TODO
	// redis 预热？？ 维持关注用户列表 redis key 的存在 维持点赞视频列表 redis key 的存在

	return &response.UserRegisterOrLogin{
		StatusCode: response.Success,
		StatusMsg:  "用户注册成功",
		Token:      token,
		UserID:     userID,
	}, nil
}

func (service *UserService) LoginService() (*response.UserRegisterOrLogin, error) {
	logTag := "service.user.Login err:"
	userDao := dao.NewUserDao()
	//先判断用户存不存在
	user, err := userDao.SelectUserByName(service.Username)
	if err != nil {
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	if user.ID == 0 {
		err := fmt.Errorf("用户不存在")
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	// TODO这里可以记录 用户输入密码错误多少次就拒绝登录了（或者用户名都不存在） redis实现

	if !util.BcryptCheck(service.Password, user.Password) {
		err := fmt.Errorf("用户密码错误")
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}

	// 签发token
	token, err := util.SignToken(user.ID)
	if err != nil {
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}

	// TODO:redis 预热？？ 这里可以预热 用户要登录 肯定立马要进入主页
	// 主页 点赞的视频列表 关注列表 粉丝  可以go协程异步执行

	return &response.UserRegisterOrLogin{
		StatusCode: response.Success,
		StatusMsg:  "用户登录成功",
		Token:      token,
		UserID:     user.ID,
	}, nil
}
