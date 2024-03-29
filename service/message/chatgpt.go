package main

import (
	"douyin/package/constant"

	"go.uber.org/zap"
)

func sendToChatGPT(userID uint64, content string) error {
	// 先将消息写入数据库
	err := CreateMessage(userID, constant.ChatGPTID, content)
	if err != nil {
		return err
	}
	go requestToChatGPT(userID, content)
	return nil
}

func requestToChatGPT(userID uint64, content string) {
	ans := requestToSparkAPI(content)
	if ans == "" {
		return
	}
	err := CreateMessage(constant.ChatGPTID, userID, ans)
	if err != nil {
		zap.L().Error(err.Error())
	}
}
