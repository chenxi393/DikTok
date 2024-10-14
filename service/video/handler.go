package main

import (
	"context"

	pbvideo "diktok/grpc/video"
	"diktok/package/util"
	"diktok/service/video/logic"

	"go.uber.org/zap"
)

type VideoService struct {
	pbvideo.UnimplementedVideoServer
}

func (s *VideoService) Feed(ctx context.Context, req *pbvideo.FeedRequest) (*pbvideo.FeedResponse, error) {
	zap.L().Sugar().Infof("[Feed] req = %s", util.GetLogStr(req))
	resp, err := logic.Feed(ctx, req)
	zap.L().Sugar().Infof("[Feed] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}

func (s *VideoService) Publish(ctx context.Context, req *pbvideo.PublishRequest) (*pbvideo.PublishResponse, error) {
	zap.L().Sugar().Infof("[Publish] req = %s", util.GetLogStr(req))
	resp, err := logic.Publish(ctx, req)
	zap.L().Sugar().Infof("[Publish] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}

func (s *VideoService) MGet(ctx context.Context, req *pbvideo.MGetReq) (*pbvideo.MGetResp, error) {
	zap.L().Sugar().Infof("[MGet] req = %s", util.GetLogStr(req))
	resp, err := logic.MGetVideos(ctx, req)
	zap.L().Sugar().Infof("[MGet] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}

func (s *VideoService) Search(ctx context.Context, req *pbvideo.SearchRequest) (*pbvideo.ListResponse, error) {
	zap.L().Sugar().Infof("[Search] req = %s", util.GetLogStr(req))
	resp, err := logic.Search(ctx, req)
	zap.L().Sugar().Infof("[Search] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}

func (s *VideoService) Pack(ctx context.Context, req *pbvideo.PackReq) (*pbvideo.PackResp, error) {
	zap.L().Sugar().Infof("[Pack] req = %s", util.GetLogStr(req))
	resp, err := logic.Pack(ctx, req)
	zap.L().Sugar().Infof("[Pack] resp = %s, err = %s", util.GetLogStr(resp), err)
	return resp, err
}
