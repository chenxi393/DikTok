package service

import (
	"douyin/database"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/llm"
	"douyin/response"
	"fmt"

	"go.uber.org/zap"
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
	// 给ChatGPT发送消息
	if service.ToUserID == llm.ChatGPTID {
		return llm.SendToChatGPT(loginUserID, service.Content)
	}
	// 发送的id是不是朋友
	isfollowing, err := cache.IsFollow(loginUserID, service.ToUserID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss)
		isfollowing, err = database.IsFollowed(loginUserID, service.ToUserID)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		go func() {
			following, err := database.SelectFollowingByUserID(loginUserID)
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			err = cache.SetFollowUserIDSet(loginUserID, following)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	if !isfollowing {
		err := fmt.Errorf("对方不是你的好友")
		return err
	}
	isfollowied, err := cache.IsFollow(service.ToUserID, loginUserID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss)
		isfollowied, err = database.IsFollowed(service.ToUserID, loginUserID)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		go func() {
			following, err := database.SelectFollowingByUserID(service.ToUserID)
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			err = cache.SetFollowUserIDSet(service.ToUserID, following)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	if !isfollowied {
		err := fmt.Errorf("对方不是你的好友")
		return err
	}
	// 发送消息考虑用不用消息队列 比较qps应该不大
	return database.CreateMessage(loginUserID, service.ToUserID, service.Content)
}

func (service *MessageService) MessageChat(loginUserID uint64) (*response.MessageResponse, error) {
	if loginUserID == service.ToUserID {
		err := fmt.Errorf("ToUserID不能是自己")
		return nil, err
	}
	// 查询所有的聊天记录 按时间顺序
	// 这里并未建立缓存
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
