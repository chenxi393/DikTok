package logic

import (
	"context"
	"errors"

	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/service/video/storage"

	"go.uber.org/zap"
)

// 企业正常 这种搜索类 列表类的接口 通过复杂查询条件是先走ES 查出IDS 再通过ID 查库
func Search(ctx context.Context, req *pbvideo.SearchRequest) (*pbvideo.ListResponse, error) {
	if req.Keyword == "" {
		return nil, errors.New(constant.BadParaRequest)
	}
	// 去数据库利用全文索引拿出所有视频数据
	videos, err := storage.SearchVideoByKeyword(req.Keyword)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 先用map 减少rpc查询次数
	userMap := make(map[int64]*pbuser.UserInfo)
	for i := range videos {
		userMap[videos[i].AuthorID] = &pbuser.UserInfo{}
	}
	videoInfo := getVideoInfo(ctx, videos, userMap, req.LoginUserId)
	return &pbvideo.ListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.SearchSuccess,
		VideoList:  videoInfo,
	}, nil

}
