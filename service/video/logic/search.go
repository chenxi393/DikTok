package logic

import (
	"context"
	"errors"

	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/service/video/storage"

	"go.uber.org/zap"
)

// 企业正常 这种搜索类 列表类的接口 通过复杂查询条件是先走ES 查出IDS 再通过ID 查库
func Search(ctx context.Context, req *pbvideo.SearchRequest) (*pbvideo.ListResponse, error) {
	if req.Keyword == "" && req.UserId == 0 {
		return nil, errors.New(constant.BadParaRequest)
	}
	var videos []*pbvideo.VideoMetaData
	if req.Keyword != "" {
		// 去数据库利用全文索引拿出所有视频数据
		videosData, err := storage.SearchVideoByKeyword(req.Keyword)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		videos = buildMGetVideosResp(videosData)
	} else if req.UserId != 0 {
		// 查询用户发布的视频
		videosData, err := MGetVideos(ctx, &pbvideo.MGetReq{
			UserId: req.UserId,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		videos = videosData.VideoList
	}
	videoInfo, err := BuildVideosInfo(ctx, nil, videos, req.GetLoginUserId())
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &pbvideo.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.SearchSuccess,
		VideoList:  videoInfo,
	}, nil

}
