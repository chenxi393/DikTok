package main

import (
	"context"
	"errors"

	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"

	"go.uber.org/zap"
)

func (s *VideoService) Search(ctx context.Context, req *pbvideo.SearchRequest) (*pbvideo.VideoListResponse, error) {
	if req.Keyword == "" {
		return nil, errors.New(constant.BadParaRequest)
	}
	// 去数据库利用全文索引拿出所有视频数据
	videos, err := SearchVideoByKeyword(req.Keyword)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 先用map 减少rpc查询次数
	userMap := make(map[int64]*pbuser.UserInfo)
	for i := range videos {
		userMap[videos[i].AuthorID] = &pbuser.UserInfo{}
	}
	videoInfo := getVideoInfo(ctx, videos, userMap, req.UserID)
	return &pbvideo.VideoListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.SearchSuccess,
		VideoList:  videoInfo,
	}, nil

}
