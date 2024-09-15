package logic

import (
	"context"

	pbmessage "diktok/grpc/message"
	"diktok/package/constant"
	"diktok/service/message/storage"

	"go.uber.org/zap"
)

func GetFirstMessage(ctx context.Context, req *pbmessage.GetFirstRequest) (*pbmessage.GetFirstResponse, error) {
	msg, err := storage.GetNewestMessage(req.UserID, req.ToUserID)
	if err != nil {
		zap.L().Error(err.Error())
	}
	// 0 表示 user send to to_user
	msgt := int32(0)
	if err != nil || msg.Content == "" {
		msg.Content = constant.DefaultMessage
	} else {
		if msg.FromUserID == req.UserID {
			msgt = 1
		}
	}
	return &pbmessage.GetFirstResponse{
		Message: msg.Content,
		MsgType: msgt,
	}, nil
}
