package logic

import (
	"context"

	pbmessage "diktok/grpc/message"
	"diktok/package/constant"
	"diktok/service/message/storage"

	"go.uber.org/zap"
)

func List(ctx context.Context, req *pbmessage.ListRequest) (*pbmessage.ListResponse, error) {
	if req.UserID == req.ToUserID {
		return &pbmessage.ListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}, nil
	}
	// 查询所有的聊天记录 按时间顺序
	// 这里并未建立缓存
	msgs, err := storage.GetMessages(req.UserID, req.ToUserID, req.PreMsgTime)
	if err != nil {
		zap.L().Error(err.Error())
		return &pbmessage.ListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.DatabaseError,
		}, nil
	}
	msgsList := make([]*pbmessage.MessageData, 0, len(msgs))
	for _, msg := range msgs {
		mm := &pbmessage.MessageData{
			Content:    msg.Content,
			CreateTime: msg.CreateTime.UnixMilli(),
			FromUserId: msg.FromUserID,
			Id:         msg.ID,
			ToUserId:   msg.ToUserID,
		}
		msgsList = append(msgsList, mm)
	}
	return &pbmessage.ListResponse{
		StatusCode:  constant.Success,
		StatusMsg:   constant.ListSuccess,
		MessageList: msgsList,
	}, nil
}
