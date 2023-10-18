package database

import (
	"douyin/model"
	"time"
)

func CreateMessage(userID, toUSerId uint64, content string) error {
	msg := model.Message{
		Content:    content,
		FromUserID: userID,
		ToUserID:   toUSerId,
		CreateTime: time.Now(),
	}
	return global_db.Model(&model.Message{}).Create(&msg).Error
}

func MessageList(userID, toUserID uint64, msgTime int64) ([]model.Message, error) {
	// 这里是客户端bug 客户端post发送评论成功后会在界面上显示 然后又请求一次数据
	// 即客户端发送消息时会显示两次重复的
	newMsgTime := time.UnixMilli(msgTime)
	msgs := make([]model.Message, 0)
	// TODO 这里用 union 避免or 不走索引的情况
	err := global_db.Raw("? UNION ? ORDER BY create_time ASC",
		global_db.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ? AND create_time > ?",
			userID, toUserID, newMsgTime),
		global_db.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ? AND create_time > ?",
			toUserID, userID, newMsgTime)).Scan(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, err
}

// 用来呈现好友列表的第一条消息
func GetMessageNewest(userID, toUserID uint64) (string, error) {
	msg := model.Message{}
	// TODO 这里用 union 避免or 不走索引的情况
	err := global_db.Raw("? UNION ? ORDER BY create_time DESC LIMIT 1",
		global_db.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ?", userID, toUserID),
		global_db.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ?", toUserID, userID)).Scan(&msg).Error
	if err != nil {
		return "", err
	}
	return msg.Content, err
}
