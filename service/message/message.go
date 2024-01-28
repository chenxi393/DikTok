package main

import (
	"context"
	pbmessage "douyin/grpc/message"
	pbrelation "douyin/grpc/relation"
	"douyin/package/constant"
	"douyin/storage/database"

	"go.uber.org/zap"
)

type MessageService struct {
	pbmessage.UnimplementedMessageServer
}

func (s *MessageService) Send(ctx context.Context, req *pbmessage.SendRequest) (*pbmessage.SendResponse, error) {
	// TODO 可能还得限制一下消息长度
	// 这里是不是给网关做？？
	if req.UserID == req.ToUserID {
		return &pbmessage.SendResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.SendToSelf,
		}, nil
	} else if req.Content == "" {
		return &pbmessage.SendResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.SendEmpty,
		}, nil
	}
	// 给ChatGPT发送消息
	if req.ToUserID == constant.ChatGPTID {
		err := sendToChatGPT(req.UserID, req.Content)
		if err != nil {
			return &pbmessage.SendResponse{
				StatusCode: constant.Failed,
				StatusMsg:  err.Error(),
			}, nil
		}
	} else {
		// 判断发送的id是不是朋友
		isFriend, err := relationClient.IsFriend(ctx, &pbrelation.ListRequest{
			LoginUserID: req.UserID,
			UserID:      req.ToUserID,
		})
		if err != nil {
			return &pbmessage.SendResponse{
				StatusCode: constant.Failed,
				StatusMsg:  err.Error(),
			}, nil
		}
		if !isFriend.Result {
			return &pbmessage.SendResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.IsNotFriend,
			}, nil
		}
		// 发送消息考虑用不用消息队列 qps应该不大
		// 后续可以用websocket完善 推送模型 TODO
		err = database.CreateMessage(req.UserID, req.ToUserID, req.Content)
		if err != nil {
			zap.L().Error(err.Error())
			return &pbmessage.SendResponse{
				StatusCode: constant.Failed,
				StatusMsg:  constant.DatabaseError,
			}, nil
		}
	}
	return &pbmessage.SendResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.SendSuccess,
	}, nil
}

func (s *MessageService) List(ctx context.Context, req *pbmessage.ListRequest) (*pbmessage.ListResponse, error) {
	if req.UserID == req.ToUserID {
		return &pbmessage.ListResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}, nil
	}
	// 查询所有的聊天记录 按时间顺序
	// 这里并未建立缓存
	msgs, err := database.MessageList(req.UserID, req.ToUserID, req.PreMsgTime)
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

func (s *MessageService) GetFirstMessage(ctx context.Context, req *pbmessage.GetFirstRequest) (*pbmessage.GetFirstResponse, error) {
	msg, err := database.GetMessageNewest(req.UserID, req.ToUserID)
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
