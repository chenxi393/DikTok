package main

import (
	"context"
	pbuser "douyin/grpc/user"
	"douyin/package/constant"
	"errors"
)

func (s *UserService) Update(ctx context.Context, req *pbuser.UpdateRequest) (*pbuser.UpdateResponse, error) {
	// TODO  TODO
	return nil, nil
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
