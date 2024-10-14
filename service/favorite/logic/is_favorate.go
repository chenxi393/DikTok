package logic

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/package/util"
	"diktok/service/favorite/storage"

	"go.uber.org/zap"
)

func IsFavorite(ctx context.Context, req *pbfavorite.IsFavoriteReq) (*pbfavorite.IsFavoriteResp, error) {
	// 获取用户的喜欢视频列表
	likingVideos, err := storage.GetFavoriteSet(req.GetUserID())
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		likingVideos, err = storage.SelectFavoriteVideoByUserID(req.GetUserID())
		if err != nil {
			return nil, err
		}
		go func() {
			err := storage.SetFavoriteSet(req.GetUserID(), likingVideos)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	mp := util.Slice2Map(likingVideos)
	resMp := make(map[int64]bool, 0)
	for _, f := range req.GetVideoID() {
		if _, ok := mp[f]; ok {
			resMp[f] = true
		}
	}
	return &pbfavorite.IsFavoriteResp{
		IsFavorite: resMp,
	}, nil
}
