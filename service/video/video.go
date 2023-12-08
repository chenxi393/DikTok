package service

import (
	"douyin/database"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/response"
	"time"

	"go.uber.org/zap"
)

type FeedService struct {
	// 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	LatestTime int64 `query:"latest_time"`
	// 用户登录状态下设置
	Token string `query:"token"`
	// 新增topic
	Topic string `query:"topic"`
}

// userID =0 表示未登录
func (service *FeedService) GetFeed(userID uint64) (*response.FeedResponse, error) {
	// TODO: 已登录可以有一个用户画像 做一个视频推荐功能
	// 直接去数据库里查出30个数据  LatestTime 限制返回视频的最晚时间
	var videos []response.VideoData
	var err error
	if service.Topic == "" {
		videos, err = database.SelectFeedVideoList(constant.MaxVideoNumber, service.LatestTime)
	} else {
		videos, err = database.SelectFeedVideoByTopic(constant.MaxVideoNumber, service.LatestTime, service.Topic)
	}
	// FIXME 这里视频数有可能为0
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	videoData := response.VideoDataInfo(videos)
	nextTime := time.Now()
	if len(videoData) != 0 {
		nextTime = videos[len(videoData)-1].PublishTime
	} else {
		//当视频数为0 的时候返回友好提示
		return &response.FeedResponse{
			StatusCode: response.Success, //返回Success 客户端下次才会更新时间
			StatusMsg:  response.NoMoreVideos,
			NextTime:   nextTime.UnixMilli(),
		}, nil
	}
	// 用户未登录直接返回
	if userID == 0 {
		return &response.FeedResponse{
			StatusCode: response.Success,
			StatusMsg:  response.FeedSuccess,
			VideoList:  videoData,
			NextTime:   nextTime.UnixMilli(),
		}, nil
	}
	// 获取用户的关注列表
	following, err := cache.GetFollowUserIDSet(userID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		following, err = database.SelectFollowingByUserID(userID)
		if err != nil {
			return nil, err
		}
		go func() {
			err := cache.SetFollowUserIDSet(userID, following)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	followingMap := make(map[uint64]struct{}, len(following))
	for _, f := range following {
		followingMap[f] = struct{}{}
	}
	// 获取用户的喜欢视频列表
	likingVideos, err := cache.GetFavoriteSet(userID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		likingVideos, err = database.SelectFavoriteVideoByUserID(userID)
		if err != nil {
			return nil, err
		}
		go func() {
			err := cache.SetFavoriteSet(userID, likingVideos)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	likingMap := make(map[uint64]struct{}, len(likingVideos))
	for _, f := range likingVideos {
		likingMap[f] = struct{}{}
	}
	// 要注意 自己的视频算被自己关注了
	// 判断是否点赞和是否关注
	followingMap[userID] = struct{}{}
	for i, rr := range videoData {
		if _, ok := followingMap[rr.Author.Id]; ok {
			videoData[i].Author.IsFollow = true
		}
		if _, ok := likingMap[rr.ID]; ok {
			videoData[i].IsFavorite = true
		}
	}

	return &response.FeedResponse{
		StatusCode: response.Success,
		StatusMsg:  response.FeedSuccess,
		VideoList:  videoData,
		NextTime:   nextTime.UnixMilli(),
	}, nil
}
