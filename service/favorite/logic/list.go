package logic

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	"diktok/package/constant"
	"diktok/service/favorite/storage"

	"go.uber.org/zap"
)

func List(ctx context.Context, req *pbfavorite.ListRequest) (*pbfavorite.ListResponse, error) {
	// TODO 加分布式锁
	// redis查找所有喜欢的视频ID
	videoIDs, err := storage.GetFavoriteSet(req.UserID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss)
		videoIDs, err = storage.SelectFavoriteVideoByUserID(req.UserID)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		// 加入到缓存里
		go func() {
			err := storage.SetFavoriteSet(req.UserID, videoIDs)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	return &pbfavorite.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FavoriteListSuccess,
		VideoList:  videoIDs,
	}, nil
}
