package main

import (
	"context"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"douyin/storage/database"
	"errors"

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
// 可以丰富用户名的判断TODO
func isUsernameOK(username string) error {
	if len(username) <= 0 || len(username) > 32 {
		return errors.New(constant.BadParaRequest)
	}
	//判断用户名存不存在
	_, err := database.SelectUserByName(username)
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
