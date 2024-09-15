package logic

import (
	"context"
	pbmessage "diktok/grpc/message"
	pbrelation "diktok/grpc/relation"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/service/message/storage"

	"go.uber.org/zap"
)

func Send(ctx context.Context, req *pbmessage.SendRequest) (*pbmessage.SendResponse, error) {
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
		isFriend, err := rpc.RelationClient.IsFriend(ctx, &pbrelation.ListRequest{
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
		err = storage.CreateMessage(req.UserID, req.ToUserID, req.Content)
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

func sendToChatGPT(userID int64, content string) error {
	// 先将消息写入数据库
	err := storage.CreateMessage(userID, constant.ChatGPTID, content)
	if err != nil {
		return err
	}
	go requestToChatGPT(userID, content)
	return nil
}

func requestToChatGPT(userID int64, content string) {
	ans := requestToSparkAPI(content)
	if ans == "" {
		return
	}
	err := storage.CreateMessage(constant.ChatGPTID, userID, ans)
	if err != nil {
		zap.L().Error(err.Error())
	}
}
