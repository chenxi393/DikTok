package logic

import (
	"context"
	"time"

	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/service/video/storage"
	"diktok/storage/database"
	"diktok/storage/database/query"

	"go.uber.org/zap"
	"gorm.io/gen"
)

func Feed(ctx context.Context, req *pbvideo.FeedRequest) (*pbvideo.FeedResponse, error) {
	// TODO: 已登录可以有一个用户画像 做一个视频推荐功能
	// TODO 目前feed 流 也没有使用缓存什么的  有粗排 细拍什么的 可以看下推荐算法
	so := query.Use(database.DB).Video
	var conds []gen.Condition
	if req.LatestTime > 0 {
		conds = append(conds, so.PublishTime.Lt(time.UnixMilli(req.GetLatestTime())))
	}
	if req.Topic != "" {
		conds = append(conds, so.Topic.Like(req.Topic+"%"))
	}

	videos, _, err := storage.MGetVideosByCond(ctx, 0, 30, conds, so.PublishTime.Desc())
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	nextTime := time.Now()
	// 这里视频数有可能为0
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].PublishTime
	} else {
		//当视频数为0 的时候返回友好提示
		return &pbvideo.FeedResponse{
			StatusCode: constant.Success, //返回Success 客户端下次才会更新时间
			StatusMsg:  constant.NoMoreVideos,
			NextTime:   nextTime.UnixMilli(),
		}, nil
	}
	// 先用map 减少rpc查询次数
	userMap := make(map[int64]*pbuser.UserInfo)
	for i := range videos {
		userMap[videos[i].AuthorID] = &pbuser.UserInfo{}
	}
	videoInfos := getVideoInfo(ctx, videos, userMap, req.LoginUserId)
	return &pbvideo.FeedResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FeedSuccess,
		VideoList:  videoInfos,
		NextTime:   nextTime.UnixMilli(),
	}, nil
}
