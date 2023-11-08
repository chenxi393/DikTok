package util

import (
	"douyin/database"
	"douyin/model"

	"go.uber.org/zap"
)

const (
	// chatgpt
	ChatGPTAvatar = "http://s2a5yl4lg.hn-bkt.clouddn.com/2022chatgpt.png"
	ChatGPTName   = "ChatGPT"
	ChatGPTID     = 1
)

func SendToChatGPT(userID uint64, content string) error {
	// 先将消息写入数据库
	err := database.CreateMessage(userID, ChatGPTID, content)
	if err != nil {
		return err
	}
	go requestToChatGPT(userID, content)
	return nil
}

func requestToChatGPT(userID uint64, content string) {
	ans := RequestToSparkAPI(content)
	if ans == "" {
		return
	}
	err := database.CreateMessage(ChatGPTID, userID, ans)
	if err != nil {
		zap.L().Error(err.Error())
	}
}

// 将chatgpt注册为用户
func RegisterChatGPT() {
	user := &model.User{
		ID:       ChatGPTID,
		Username: ChatGPTName,
		Avatar:   ChatGPTAvatar,
	}
	_, err := database.CreateUser(user)
	if err != nil {
		zap.Error(err)
	}
}
