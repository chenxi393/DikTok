package service

import (
	"douyin/database"
	"douyin/model"
	"douyin/response"
	"fmt"

	"go.uber.org/zap"
)

type FavoriteService struct {
	// 1-点赞，2-取消点赞
	ActionType string `json:"action_type"`
	// 用户鉴权token
	Token string `json:"token"`
	// 视频id
	VideoID uint64 `json:"video_id"`
	// 要查询的用户id
	UserID uint64 `json:"user_id"`
}

func (service *FavoriteService) FavoriteAction(userID uint64) error {
	// TODO 这里用到了消息队列
	return database.FavoriteVideo(userID, service.VideoID)
}

func (service *FavoriteService) UnFavoriteAction(userID uint64) error {
	// TODO 这里用到了消息队列
	return database.UnFavoriteVideo(userID, service.VideoID)
}

func (service *FavoriteService) FavoriteList(userID uint64) ([]response.Video, error) {
	// 先查找所有喜欢的视频ID
	videoIDs, err := database.SelectFavoriteVideoByUserID(userID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 然后去数据库批量查找视频数据
	videos, err := database.SelectVideoListByVideoID(videoIDs)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 拿到视频数据之后 还得一个视频一个视频拿到作者信息
	userIDs := make([]uint64, 0, len(videoIDs))
	for _, video := range videos {
		userIDs = append(userIDs, video.AuthorID)
	}
	// 批量拿到作者信息 但是还需要填空 哪个作者对应哪个
	usersData, err := database.SelectUserListByIDs(userIDs)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 判断用户有没有关注 获取用户关注列表
	following, err := database.SelectFollowingByUserID(userID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	followingMap := make(map[uint64]struct{}, len(following))
	for _, i := range following {
		followingMap[i] = struct{}{}
	}
	usersMap := make(map[uint64]*model.User, len(usersData))
	for i, id := range usersData {
		usersMap[id.ID] = &usersData[i]
	}
	videoResponse := make([]response.Video, 0, len(videos))
	for _, video := range videos {
		if _, ok := usersMap[video.AuthorID]; ok {
			_, isFollowing := followingMap[video.AuthorID]
			vv := response.Video{
				Author:        *response.UserInfo(usersMap[video.AuthorID], isFollowing),
				CommentCount:  video.CommentCount,
				CoverURL:      video.CoverURL,
				FavoriteCount: video.FavoriteCount,
				ID:            video.ID,
				IsFavorite:    true,
				PlayURL:       video.PlayURL,
				Title:         video.Title,
			}
			videoResponse = append(videoResponse, vv)
		} else {
			err := fmt.Errorf("视频缺少作者 服务端bug")
			zap.L().Error(err.Error())
			return nil, err
		}
	}
	return videoResponse, nil
}
