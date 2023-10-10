package service

import (
	"douyin/dal/dao"
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
	logTag := "service.feed.GetFeed err:"
	videoData, err := dao.SelectFeedVideoList(maxVideoNum, service.LatestTime)
	if err != nil {
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	var nextTime time.Time
	if len(videoData) != 0 {
		nextTime, err = dao.SelectPublishTimeByVideoID(videoData[len(videoData)-1].ID)
		if err != nil {
			zap.L().Error(logTag, zap.Error(err))
			return nil, err
		}
	}
	// 查看有没有token
	if service.Token == nil || *service.Token == "" {
		return &response.FeedResponse{
			StatusCode: response.Success,
			StatusMsg:  response.FeedSuccess,
			VideoList:  videoData,
			NextTime:   nextTime.UnixMilli(),
		}, nil
	}
	userClaim, err := util.ParseToken(*service.Token)
	if err != nil || userClaim.UserID == 0 {
		err := fmt.Errorf("解析token失败")
		zap.L().Error(logTag, zap.Error(err))
		return nil, err
	}
	// 获取用户的关注列表
	// TODO 去redis拿关注列表
	following, err := dao.SelectFollowingByUserID(userClaim.UserID)
	if err != nil {
		return nil, err
	}
	followingMap := make(map[uint64]struct{}, len(following))
	for _, f := range following {
		followingMap[f] = struct{}{}
	}
	// 获取用户的喜欢视频列表
	likingVideos, err := dao.SelectFavoriteVideoByUserID(userClaim.UserID)
	if err != nil {
		return nil, err
	}
	likingMap := make(map[uint64]struct{}, len(likingVideos))
	for _, f := range likingVideos {
		likingMap[f] = struct{}{}
	}
	// 判断是否点赞和是否关注
	for _, rr := range videoData {
		if _, ok := followingMap[rr.ID]; ok {
			rr.IsFollow = true
		}
		if _, ok := likingMap[rr.ID]; ok {
			rr.IsFavorite = true
		}
	}

	return &response.FeedResponse{
		StatusCode: response.Success,
		VideoList:  videoData,
		NextTime:   nextTime.UnixMilli(),
	}, nil
}
