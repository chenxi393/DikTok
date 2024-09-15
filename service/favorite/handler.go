package main

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	"diktok/package/util"
	"diktok/service/favorite/logic"

	"go.uber.org/zap"
)

type FavoriteService struct {
	pbfavorite.UnimplementedFavoriteServer
}

func (s *FavoriteService) Like(ctx context.Context, req *pbfavorite.LikeRequest) (*pbfavorite.LikeResponse, error) {
	zap.L().Sugar().Infof("[Like] req = %s", util.GetLogStr(req))
	resp, err := logic.Like(ctx, req)
	zap.L().Sugar().Infof("[Like] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *FavoriteService) Unlike(ctx context.Context, req *pbfavorite.LikeRequest) (*pbfavorite.LikeResponse, error) {
	zap.L().Sugar().Infof("[Unlike] req = %s", util.GetLogStr(req))
	resp, err := logic.Unlike(ctx, req)
	zap.L().Sugar().Infof("[Unlike] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *FavoriteService) List(ctx context.Context, req *pbfavorite.ListRequest) (*pbfavorite.ListResponse, error) {
	zap.L().Sugar().Infof("[List] req = %s", util.GetLogStr(req))
	resp, err := logic.List(ctx, req)
	zap.L().Sugar().Infof("[List] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *FavoriteService) IsFavorite(ctx context.Context, req *pbfavorite.IsFavoriteRequest) (*pbfavorite.IsFavoriteResponse, error) {
	zap.L().Sugar().Infof("[IsFavorite] req = %s", util.GetLogStr(req))
	resp, err := logic.IsFavorite(ctx, req)
	zap.L().Sugar().Infof("[IsFavorite] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *FavoriteService) Count(ctx context.Context, req *pbfavorite.CountReq) (*pbfavorite.CountResp, error) {
	zap.L().Sugar().Infof("[Count] req = %s", util.GetLogStr(req))
	resp, err := logic.Count(ctx, req)
	zap.L().Sugar().Infof("[Count] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}
