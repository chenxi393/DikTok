package main

import (
	"context"

	pbfavorite "diktok/grpc/favorite"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"

	"go.uber.org/zap"
)

type FavoriteService struct {
	pbfavorite.UnimplementedFavoriteServer
}

func (s *FavoriteService) Like(ctx context.Context, req *pbfavorite.LikeRequest) (*pbfavorite.LikeResponse, error) {
	// TODO 可以拿redis限制一下用户点赞的速率 比如1分钟只能点赞10次
	err := FavoriteVideo(req.UserID, req.VideoID, 1)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}

	return &pbfavorite.LikeResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FavoriteSuccess,
	}, nil
}

func (s *FavoriteService) Unlike(ctx context.Context, req *pbfavorite.LikeRequest) (*pbfavorite.LikeResponse, error) {
	err := FavoriteVideo(req.UserID, req.VideoID, -1)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}
	return &pbfavorite.LikeResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.UnFavoriteSuccess,
	}, nil
}

func (s *FavoriteService) List(ctx context.Context, req *pbfavorite.ListRequest) (*pbfavorite.ListResponse, error) {
	// TODO 加分布式锁
	// redis查找所有喜欢的视频ID
	videoIDs, err := GetFavoriteSet(req.UserID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss)
		videoIDs, err = SelectFavoriteVideoByUserID(req.UserID)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		// 加入到缓存里
		go func() {
			err := SetFavoriteSet(req.UserID, videoIDs)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	// FIXME 这里非顺序返回
	// 返回的是按id倒叙返回
	videos, err := videoClient.GetVideosByUserID(ctx, &pbvideo.GetVideosRequest{
		UserID:  req.LoginUserID,
		VideoID: videoIDs,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &pbfavorite.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FavoriteListSuccess,
		VideoList:  videos.GetVideoList(),
	}, nil
}

func (s *FavoriteService) IsFavorite(ctx context.Context, req *pbfavorite.LikeRequest) (*pbfavorite.IsFavoriteResponse, error) {
	// 获取用户的喜欢视频列表
	likingVideos, err := GetFavoriteSet(req.GetUserID())
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		likingVideos, err = SelectFavoriteVideoByUserID(req.GetUserID())
		if err != nil {
			return nil, err
		}
		go func() {
			err := SetFavoriteSet(req.GetUserID(), likingVideos)
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
