package service

import (
	"douyin/database"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/util"
	"douyin/response"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type FeedService struct {
	// 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	LatestTime *int64 `query:"latest_time"`
	// 用户登录状态下设置
	Token *string `query:"token"`
}

var maxVideoNum = 30

func (service *FeedService) GetFeed() (*response.FeedResponse, error) {
	// TODO: 已登录可以有一个用户画像 做一个视频推荐功能
	// 直接去数据库里查出30个数据  LatestTime 限制返回视频的最晚时间
	videos, err := database.SelectFeedVideoList(maxVideoNum, service.LatestTime)
	// FIXME 这里视频数有可能为0
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	videoData := response.VideoDataInfo(videos)
	nextTime := time.Now()
	if len(videoData) != 0 {
		nextTime, err = database.SelectPublishTimeByVideoID(videoData[len(videoData)-1].ID)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
	} else { //当视频数为0 的时候返回友好提示
		return &response.FeedResponse{
			StatusCode: response.Success, //返回成功 客户端下次才会更新时间
			StatusMsg:  "视频见底了",
			VideoList:  nil,
			NextTime:   nextTime.UnixMilli(),
		}, nil
	}
	// 查看有没有token  视频数为0也立刻返回
	if service.Token == nil || *service.Token == "" {
		return &response.FeedResponse{
			StatusCode: response.Success,
			StatusMsg:  response.FeedSuccess,
			VideoList:  videoData,
			NextTime:   nextTime.UnixMilli(),
		}, nil
	}
	claim, err := util.ParseToken(*service.Token)
	if err != nil || claim.UserID == 0 {
		err := fmt.Errorf("解析token失败 请重新登录") // 这里哪怕鉴权失败页给用户返回信息
		zap.L().Error(err.Error())
		return nil, err
	}
	// 获取用户的关注列表
	following, err := cache.GetFollowUserIDSet(claim.UserID)
	if err != nil {
		zap.L().Error(constant.CacheMiss)
		following, err = database.SelectFollowingByUserID(claim.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			err := cache.SetFollowUserIDSet(claim.UserID, following)
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
	likingVideos, err := cache.GetFavoriteSet(claim.UserID)
	if err != nil {
		zap.L().Error(constant.CacheMiss)
		likingVideos, err = database.SelectFavoriteVideoByUserID(claim.UserID)
		if err != nil {
			return nil, err
		}
		go func() {
			err := cache.SetFavoriteSet(claim.UserID, likingVideos)
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
	followingMap[claim.UserID] = struct{}{}
	for i, rr := range videoData {
		if _, ok := followingMap[rr.Author.ID]; ok {
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
