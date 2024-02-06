package main

import (
	"context"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"douyin/storage/database"
	"errors"
	"regexp"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	updateUsername   = 1
	updatePassword   = 2
	updateSignature  = 3
	updateAvatar     = 4
	updateBackground = 5
)

func (s *UserService) Update(ctx context.Context, req *pbuser.UpdateRequest) (*pbuser.UpdateResponse, error) {
	user, err := database.SelectUserByID(req.UserID)
	if err != nil {
		otelzap.L().Error(constant.DatabaseError, zap.Error(err))
		return nil, errors.New(constant.UserNoExist)
	}
	switch req.UpdateType {
	case updateUsername:
		// 旧密码判断
		if !bcryptCheck(req.OldPassword, user.Password) {
			return nil, errors.New(constant.SecretError)
		}
		// 用户名合法判断
		err := isUsernameOK(req.Username)
		if err != nil {
			return nil, err
		}
		user.Username = req.Username
	case updatePassword:
		// 旧密码判断
		if !bcryptCheck(req.OldPassword, user.Password) {
			return nil, errors.New(constant.SecretError)
		}
		err := isPasswordOK(req.NewPassword)
		if err != nil {
			return nil, err
		}
		encryptedPassword := bcryptHash(req.NewPassword)
		user.Password = encryptedPassword
	case updateSignature:
		user.Signature = req.Signature
	case updateAvatar, updateBackground:
		// TODO 上传图片怎么做
		return nil, nil
	}
	err = database.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return &pbuser.UpdateResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.UpdateSuccess,
	}, nil
}

func isPasswordOK(password string) error {
	if len(password) < 6 || len(password) > 32 {
		return errors.New(constant.SecretFormatError)
	}
	// TODO 复杂度判断 可以使用正则 记得去除常数
	if password == constant.EasySecret {
		return errors.New(constant.SecretFormatEasy)
	}
	return nil
}

// 判断用户名是否合法 是否已经存在该用户名
func isUsernameOK(username string) error {
	// 正则表达式判断 4-12个字符 英文 数字 下划线
	/*
		^ 表示字符串的开始
		[a-zA-Z0-9_] 表示允许使用字母、数字和下划线
		{4,12} 表示字符数限制在4到12个之间
		$ 表示字符串的结束
	*/
	regex := `^[a-zA-Z0-9_]{4,12}$`
	match, err := regexp.MatchString(regex, username)
	if err != nil {
		return err
	}
	if !match {
		return errors.New(constant.UsernameFormatErr)
	}

	//判断用户名存不存在
	_, err = database.SelectUserByName(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error(constant.DatabaseError, zap.Error(err))
		return err
	}
	// 这一块逻辑有点奇怪
	if err != gorm.ErrRecordNotFound {
		return errors.New(constant.UserDepulicate)
	}
	return nil
}
