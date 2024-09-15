package main

import (
	"context"

	pbuser "diktok/grpc/user"
	"diktok/package/util"
	"diktok/service/user/logic"

	"go.uber.org/zap"
)

type UserService struct {
	pbuser.UnimplementedUserServer
}

func (s *UserService) Register(ctx context.Context, req *pbuser.RegisterRequest) (*pbuser.RegisterResponse, error) {
	zap.L().Sugar().Infof("[Register] req = %s", util.GetLogStr(req))
	resp, err := logic.Register(ctx, req)
	zap.L().Sugar().Infof("[Register] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *UserService) Login(ctx context.Context, req *pbuser.LoginRequest) (*pbuser.LoginResponse, error) {
	zap.L().Sugar().Infof("[Login] req = %s", util.GetLogStr(req))
	resp, err := logic.Login(ctx, req)
	zap.L().Sugar().Infof("[Login] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *UserService) Info(ctx context.Context, req *pbuser.InfoRequest) (*pbuser.InfoResponse, error) {
	zap.L().Sugar().Infof("[Info] req = %s", util.GetLogStr(req))
	resp, err := logic.Info(ctx, req)
	zap.L().Sugar().Infof("[Info] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *UserService) List(ctx context.Context, req *pbuser.ListReq) (*pbuser.ListResp, error) {
	zap.L().Sugar().Infof("[List] req = %s", util.GetLogStr(req))
	resp, err := logic.List(ctx, req)
	zap.L().Sugar().Infof("[List] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *UserService) Update(ctx context.Context, req *pbuser.UpdateRequest) (*pbuser.UpdateResponse, error) {
	zap.L().Sugar().Infof("[Update] req = %s", util.GetLogStr(req))
	resp, err := logic.Update(ctx, req)
	zap.L().Sugar().Infof("[Update] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}
