package service

import (
	"douyin/database"
	"douyin/response"
	"fmt"
)

type MessageService struct {
	// 1-发送消息
	ActionType string `query:"action_type"`
	// 消息内容
	Content string `query:"content"`
	// 对方用户id
	ToUserID uint64 `query:"to_user_id"`
	// 用户鉴权token
	Token string `query:"token"`
	// //上次最新消息的时间（新增字段-apk更新中）
	Pre_msg_time int64 `query:"pre_msg_time"`
}

func (service *MessageService) MessageAction(loginUserID uint64) error {
	// TODO 可能还得限制一下消息长度
	if loginUserID == service.ToUserID {
		err := fmt.Errorf("不能给自己发送消息")
		return err
	} else if service.Content == "" {
		err := fmt.Errorf("消息内容为空 发送失败")
		return err
	} else if service.ActionType != "1" {
		err := fmt.Errorf("ActionType 错误")
		return err
	}
	return database.CreateMessage(loginUserID, service.ToUserID, service.Content)
}

func (service *MessageService) MessageChat(loginUserID uint64) (*response.MessageResponse, error) {
	if loginUserID == service.ToUserID {
		err := fmt.Errorf("ToUserID不能是自己")
		return nil, err
	}
	// 查询所有的聊天记录 按时间顺序
	msgs, err := database.MessageList(loginUserID, service.ToUserID, service.Pre_msg_time)
	if err != nil {
		return nil, err
	}
	msgsList := make([]response.Message, 0, len(msgs))
	for _, msg := range msgs {
		mm := response.Message{
			Content:    msg.Content,
			CreateTime: msg.CreateTime.UnixMilli(),
			FromUserID: msg.FromUserID,
			ID:         msg.ID,
			ToUserID:   msg.ToUserID,
		}
		msgsList = append(msgsList, mm)
	}
	return &response.MessageResponse{
		MessageList: msgsList,
		StatusCode:  response.Success,
		StatusMsg:   "消息列表加载成功",
	}, nil
}
