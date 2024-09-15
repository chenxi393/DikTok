package main

import (
	"context"

	pbcomment "diktok/grpc/comment"
	"diktok/package/util"
	"diktok/service/comment/logic"

	"go.uber.org/zap"
)

type CommentService struct {
	pbcomment.UnimplementedCommentServer
}

func (s *CommentService) Add(ctx context.Context, req *pbcomment.AddRequest) (*pbcomment.CommentResponse, error) {
	zap.L().Sugar().Infof("[Add] req = %s", util.GetLogStr(req))
	resp, err := logic.Add(ctx, req)
	zap.L().Sugar().Infof("[Add] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *CommentService) Delete(ctx context.Context, req *pbcomment.DeleteRequest) (*pbcomment.CommentResponse, error) {
	zap.L().Sugar().Infof("[Delete] req = %s", util.GetLogStr(req))
	resp, err := logic.Delete(ctx, req)
	zap.L().Sugar().Infof("[Delete] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *CommentService) List(ctx context.Context, req *pbcomment.ListRequest) (*pbcomment.ListResponse, error) {
	zap.L().Sugar().Infof("[List] req = %s", util.GetLogStr(req))
	resp, err := logic.List(ctx, req)
	zap.L().Sugar().Infof("[List] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *CommentService) Count(ctx context.Context, req *pbcomment.CountReq) (*pbcomment.CountResp, error) {
	zap.L().Sugar().Infof("[Count] req = %s", util.GetLogStr(req))
	resp, err := logic.Count(ctx, req)
	zap.L().Sugar().Infof("[Count] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}
