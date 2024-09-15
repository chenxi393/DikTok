package logic

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	"diktok/service/favorite/storage"

	"go.uber.org/zap"
)

func Count(ctx context.Context, req *pbfavorite.CountReq) (*pbfavorite.CountResp, error) {
	total, err := storage.CountmByVideoIDs(ctx, req.GetVideoID())
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &pbfavorite.CountResp{
		Total: total,
	}, nil
}
