package main

import (
	"context"
	"douyin/config"
	pbfavorite "douyin/grpc/favorite"
	pbuser "douyin/grpc/user"
	pbvideo "douyin/grpc/video"
	"douyin/model"
	"douyin/package/constant"
	"sync"
	"time"

	"go.uber.org/zap"
)

type VideoService struct {
	pbvideo.UnimplementedVideoServer
}

// userID =0 表示未登录
func (s *VideoService) Feed(ctx context.Context, req *pbvideo.FeedRequest) (*pbvideo.FeedResponse, error) {
	zap.L().Sugar().Infof("%+v", req)
	// TODO: 已登录可以有一个用户画像 做一个视频推荐功能
	// 直接去数据库里查出30个数据  LatestTime 限制返回视频的最晚时间
	var videos []*model.Video
	var err error
	if req.Topic == "" {
		videos, err = SelectFeedVideoList(constant.MaxVideoNumber, req.LatestTime)
	} else {
		videos, err = SelectFeedVideoByTopic(constant.MaxVideoNumber, req.LatestTime, req.Topic)
	}
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
	videoInfos := getVideoInfo(ctx, videos, userMap, req.UserID)
	return &pbvideo.FeedResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.FeedSuccess,
		VideoList:  videoInfos,
		NextTime:   nextTime.UnixMilli(),
	}, nil
}

func (s *VideoService) GetVideosByUserID(ctx context.Context, req *pbvideo.GetVideosRequest) (*pbvideo.GetVideosResponse, error) {
	// 去数据库批量查找视频数据
	// TODO 要去redis查找视频信息 否则存视频没有意义
	// 但是我又觉得最终还是要走DB 一开始走不就行
	// 再看看怎么写合理 暂时只走数据库（db肯定是顺序的）
	// 最重点的redis返回的videoIDs 不是顺序的
	// 那么走redis查到的数据是乱序的（用zset解决 但是代码复杂）
	videos, err := SelectVideosByVideoID(req.VideoID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 先用map 减少rpc查询次数
	userMap := make(map[int64]*pbuser.UserInfo)
	for i := range videos {
		userMap[videos[i].AuthorID] = &pbuser.UserInfo{}
	}
	videoInfos := getVideoInfo(ctx, videos, userMap, req.UserID)
	return &pbvideo.GetVideosResponse{
		VideoList: videoInfos,
	}, nil

}
func videoDataInfo(v *model.Video, u *pbuser.UserInfo) *pbvideo.VideoData {
	return &pbvideo.VideoData{
		Author:        u,
		Id:            v.ID,
		PlayUrl:       config.System.Qiniu.OssDomain + "/" + v.PlayURL,
		CoverUrl:      config.System.Qiniu.OssDomain + "/" + v.CoverURL,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		Title:         v.Title,
		Topic:         v.Topic,
		PublishTime:   v.PublishTime.Format("2006-01-02 15:04"),
	}
}

// RPC调用拿userMap 里的用户信息 拿video里的详细信息 返回
func getVideoInfo(ctx context.Context, videos []*model.Video, userMap map[int64]*pbuser.UserInfo, loginUserID int64) []*pbvideo.VideoData {
	// rpc调用 去拿个人信息
	wg := &sync.WaitGroup{}
	wg.Add(len(userMap))
	for userID := range userMap {
		go func(id int64) {
			defer wg.Done()
			// TODO 这里是不是也应该 rpc批量拿出来 而不是一个个去拿
			user, err := userClient.Info(ctx, &pbuser.InfoRequest{
				LoginUserID: loginUserID,
				UserID:      id,
			})
			if err != nil {
				zap.L().Error(err.Error())
			}
			if err == nil && user.StatusCode != 0 {
				zap.L().Error("rpc 调用错误")
			}
			// 这里map会不会有并发问题啊
			// TODO 去测试一下
			// 这里如果不用 指针写入的化 会导致下面 videoInfo
			// append 地址被改变 要不就上锁 所有rpc请求之后 再下一个
			// 但是这里之间 直接使用* 似乎也不太好 报了warning
			// 说内部有锁  不能复制 TODO
			*userMap[id] = *user.GetUser()
		}(userID)
	}
	videoInfos := make([]*pbvideo.VideoData, 0, len(videos))
	for i := range videos {
		videoInfos = append(videoInfos, videoDataInfo(videos[i], userMap[videos[i].AuthorID]))
	}
	// 判断请求用户是否喜欢
	if loginUserID != 0 {
		wg.Add(len(videoInfos))
		for i := range videoInfos {
			go func(i int) {
				defer wg.Done()
				result, err := favoriteClient.IsFavorite(ctx, &pbfavorite.LikeRequest{
					UserID:  loginUserID,
					VideoID: videoInfos[i].Id,
				})
				if err != nil {
					zap.L().Error(err.Error())
					return
				}
				videoInfos[i].IsFavorite = result.IsFavorite
			}(i)
		}
	}
	wg.Wait()
	return videoInfos
}
