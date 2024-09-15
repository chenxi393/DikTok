package logic

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/service/favorite/storage"

	"go.uber.org/zap"
)

func Like(ctx context.Context, req *pbfavorite.LikeRequest) (*pbfavorite.LikeResponse, error) {
	// TODO 可以拿redis限制一下用户点赞的速率 比如1分钟只能点赞10次
	err := storage.FavoriteVideo(req.UserID, req.VideoID, 1)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}

	return &pbfavorite.LikeResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FavoriteSuccess,
	}, nil
}
