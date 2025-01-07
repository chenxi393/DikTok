package response

import (
	pbmessage "diktok/grpc/message"
	"diktok/package/constant"
)

type MessageResponse struct {
	// 用户列表
	MessageList []*Message `json:"message_list"`
	// 状态码，0-成功，其他值-失败
	StatusCode int `json:"status_code"`
	// 返回状态描述
	StatusMsg string `json:"status_msg"`
}

// Message
type Message struct {
	// 消息内容
	Content string `json:"content"`
	// 消息发送时间 yyyy-MM-dd HH:MM:ss
	CreateTime int64 `json:"create_time"`
	// 消息发送者id
	FromUserID int64 `json:"from_user_id"`
	// 消息id
	ID int64 `json:"id"`
	// 消息接收者id
	ToUserID int64 `json:"to_user_id"`
}

func BuildMessageRes(msgRes *pbmessage.ListResponse) *MessageResponse {
	msgData := msgRes.MessageList
	res := &MessageResponse{
		MessageList: make([]*Message, 0, len(msgData)),
		StatusCode:  constant.Success,
		StatusMsg:   constant.FeedSuccess,
	}
	for _, v := range msgData {
		if v != nil {
			res.MessageList = append(res.MessageList, BuildMessage(v))
		}
	}
	return res
}

func BuildMessage(item *pbmessage.MessageData) *Message {
	return &Message{
		ID:         item.Id,
		Content:    item.Content,
		CreateTime: item.CreateTime,
		FromUserID: item.FromUserId,
		ToUserID:   item.ToUserId,
	}
}
