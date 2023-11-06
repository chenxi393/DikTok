package util

import (
	"douyin/config"
	"douyin/database"
	"douyin/model"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const (
	// chatgpt
	ChatGPTAvatar = "http://s2a5yl4lg.hn-bkt.clouddn.com/2022chatgpt.png"
	ChatGPTName   = "ChatGPT"
	ChatGPTID     = 1
)

type ChatGPTReply struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Created int    `json:"created"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Object  string `json:"object"`
	Choices []struct {
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func SendToChatGPT(userID uint64, content string) error {
	// 先将消息写入数据库
	err := database.CreateMessage(userID, ChatGPTID, content)
	if err != nil {
		return err
	}
	go func() {
		url := "https://api.perplexity.ai/chat/completions"
		payload := strings.NewReader("{\"model\":\"pplx-70b-chat-alpha\",\"messages\":[{\"role\":\"system\",\"content\":\"全部使用中文回复\"},{\"role\":\"user\",\"content\":\"" + content + "\"}]}")
		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "Bearer "+config.System.GPTSecret)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		var replyJSON ChatGPTReply
		err = json.Unmarshal(body, &replyJSON)
		if err != nil {
			zap.L().Error(err.Error())
		}
		if len(replyJSON.Choices) == 0 {
			zap.L().Error("大模型未回复")
			return
		}
		err = database.CreateMessage(ChatGPTID, userID, replyJSON.Choices[0].Message.Content)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}()
	return nil
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
