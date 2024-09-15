package logic

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/service/favorite/storage"

	"go.uber.org/zap"
)

func Unlike(ctx context.Context, req *pbfavorite.LikeRequest) (*pbfavorite.LikeResponse, error) {
	err := storage.FavoriteVideo(req.UserID, req.VideoID, -1)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}
	return &pbfavorite.LikeResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.UnFavoriteSuccess,
	}, nil
}
