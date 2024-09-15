package logic

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/service/favorite/storage"

	"go.uber.org/zap"
)

func IsFavorite(ctx context.Context, req *pbfavorite.IsFavoriteRequest) (*pbfavorite.IsFavoriteResponse, error) {
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
	for _, f := range likingVideos {
		if f == req.VideoID {
			return &pbfavorite.IsFavoriteResponse{
				IsFavorite: true,
			}, nil
		}
	}
	return &pbfavorite.IsFavoriteResponse{
		IsFavorite: false,
	}, nil
}
