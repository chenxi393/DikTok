package storage

import (
	"time"

	"diktok/storage/database"
	"diktok/storage/database/model"
)

func CreateMessage(userID, toUSerId int64, content string) error {
	msg := model.Message{
		Content:    content,
		FromUserID: userID,
		ToUserID:   toUSerId,
		CreateTime: time.Now(),
	}
	return database.DB.Model(&model.Message{}).Create(&msg).Error
}

func GetMessages(userID, toUserID int64, msgTime int64) ([]model.Message, error) {
	// 这里是客户端bug 客户端post发送评论成功后会在界面上显示 然后又请求一次数据
	// 即客户端发送消息时会显示两次重复的
	newMsgTime := time.UnixMilli(msgTime)
	msgs := make([]model.Message, 0)
	// 这里用 union 避免or 不走索引的情况 or两侧必须都走索引 括号也没用
	err := database.DB.Raw("? UNION ? ORDER BY create_time ASC",
		database.DB.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ? AND create_time > ?",
			userID, toUserID, newMsgTime),
		database.DB.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ? AND create_time > ?",
			toUserID, userID, newMsgTime)).Scan(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, err
}

// 用来呈现好友列表的第一条消息
func GetNewestMessage(userID, toUserID int64) (model.Message, error) {
	msg := model.Message{}
	// 这里用 union 避免or 不走索引的情况 or两侧必须都走索引 括号也没用
	err := database.DB.Raw("? UNION ? ORDER BY create_time DESC LIMIT 1",
		database.DB.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ?", userID, toUserID),
		database.DB.Model(&model.Message{}).Where("from_user_id = ? AND to_user_id = ?", toUserID, userID)).Scan(&msg).Error
	if err != nil {
		return msg, err
	}
	return msg, err
}
