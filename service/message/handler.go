package main

import (
	"context"

	pbmessage "diktok/grpc/message"
	"diktok/package/util"
	"diktok/service/message/logic"

	"go.uber.org/zap"
)

type MessageService struct {
	pbmessage.UnimplementedMessageServer
}

func (s *MessageService) Send(ctx context.Context, req *pbmessage.SendRequest) (*pbmessage.SendResponse, error) {
	zap.L().Sugar().Infof("[Send] req = %s", util.GetLogStr(req))
	resp, err := logic.Send(ctx, req)
	zap.L().Sugar().Infof("[Send] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}

func (s *MessageService) List(ctx context.Context, req *pbmessage.ListRequest) (*pbmessage.ListResponse, error) {
	zap.L().Sugar().Infof("[List] req = %s", util.GetLogStr(req))
	resp, err := logic.List(ctx, req)
	zap.L().Sugar().Infof("[List] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}

func (s *MessageService) GetFirstMessage(ctx context.Context, req *pbmessage.GetFirstRequest) (*pbmessage.GetFirstResponse, error) {
	zap.L().Sugar().Infof("[GetFirstMessage] req = %s", util.GetLogStr(req))
	resp, err := logic.GetFirstMessage(ctx, req)
	zap.L().Sugar().Infof("[GetFirstMessage] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}
