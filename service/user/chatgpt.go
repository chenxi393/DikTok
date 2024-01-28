package main

import (
	"douyin/model"
	"douyin/package/constant"
	"douyin/storage/database"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// 将chatgpt注册为用户
func registerChatGPT() {
	user := &model.User{
		ID:       constant.ChatGPTID,
		Username: constant.ChatGPTName,
		Avatar:   constant.ChatGPTAvatar,
	}
	_, err := database.CreateUser(user)
	if err != nil {
		otelzap.L().Info("ChatGPT已写入user表")
	}
}
