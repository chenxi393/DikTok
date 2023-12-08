package service

import (
	"douyin/config"
	"douyin/database"
	"douyin/model"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/mq"
	"douyin/response"
	"errors"

	"go.uber.org/zap"
)

type FavoriteService struct {
	// 1-点赞，2-取消点赞
	ActionType string `query:"action_type"`
	// 用户鉴权token
	Token string `query:"token"`
	// 视频id
	VideoID uint64 `query:"video_id"`
	// 要查询的用户id
	UserID uint64 `query:"user_id"`
}

func (service *FavoriteService) Favorite(userID uint64) (*response.CommonResponse, error) {
	// TODO 可以拿redis限制一下用户点赞的速率 比如1分钟只能点赞10次
	err := mq.SendFavoriteMessage(userID, service.VideoID, 1)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  constant.FavoriteSuccess,
	}, nil
}

func (service *FavoriteService) UnFavorite(userID uint64) (*response.CommonResponse, error) {
	err := mq.SendFavoriteMessage(userID, service.VideoID, -1)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  constant.UnFavoriteSuccess,
	}, nil
}

func (service *FavoriteService) FavoriteList(userID uint64) ([]response.Video, error) {
	// TODO 加分布式锁
	// redis查找所有喜欢的视频ID
	videoIDs, err := cache.GetFavoriteSet(service.UserID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss)
		videoIDs, err = database.SelectFavoriteVideoByUserID(service.UserID)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		// 加入到缓存里
		go func() {
			err := cache.SetFavoriteSet(service.UserID, videoIDs)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	// 然后去数据库批量查找视频数据
	// TODO 要去redis查找视频信息 否则存视频没有意义
	// 但是我又觉得最终还是要走DB 一开始走不就行
	// 再看看怎么写合理 暂时只走数据库（db肯定是顺序的）
	// 最重点的redis返回的videoIDs 不是顺序的
	// 那么走redis查到的数据是乱序的（用zset解决 但是代码复杂）
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
	// TODO 作者信息也应该先去redis里面拿 不然没有意义
	usersData, err := database.SelectUserListByIDs(userIDs)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 判断用户有没有关注 获取用户关注列表
	var following []uint64
	if userID != 0 {
		following, err = cache.GetFollowUserIDSet(userID)
		if err != nil {
			zap.L().Warn(constant.CacheMiss)
			following, err = database.SelectFollowingByUserID(userID)
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
			go func() {
				err := cache.SetFollowUserIDSet(userID, following)
				if err != nil {
					zap.L().Error(err.Error())
				}
			}()
		}
	}
	followingMap := make(map[uint64]struct{}, len(following))
	for _, i := range following {
		followingMap[i] = struct{}{}
	}
	// 把自己放进去
	followingMap[userID] = struct{}{}
	usersMap := make(map[uint64]*model.User, len(usersData))
	for i, id := range usersData {
		usersMap[id.ID] = &usersData[i]
	}
	var favorite []uint64
	if userID != 0 {
		// 获取登录用户点赞列表
		favorite, err = cache.GetFavoriteSet(userID)
		if err != nil {
			zap.L().Warn(constant.CacheMiss)
			favorite, err = database.SelectFavoriteVideoByUserID(userID)
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
			go func() {
				err := cache.SetFavoriteSet(userID, favorite)
				if err != nil {
					zap.L().Error(err.Error())
				}
			}()
		}
	}
	favoriteMap := make(map[uint64]struct{}, len(favorite))
	for _, i := range favorite {
		favoriteMap[i] = struct{}{}
	}
	videoResponse := make([]response.Video, 0, len(videos))
	for _, video := range videos {
		if _, ok := usersMap[video.AuthorID]; ok {
			_, isFollowing := followingMap[video.AuthorID]
			vv := response.Video{
				Author:        *response.UserInfo(usersMap[video.AuthorID], isFollowing),
				CommentCount:  video.CommentCount,
				CoverURL:      config.System.Qiniu.OssDomain + "/" + video.CoverURL,
				FavoriteCount: video.FavoriteCount,
				ID:            video.ID,
				IsFavorite:    false,
				PlayURL:       config.System.Qiniu.OssDomain + "/" + video.PlayURL,
				Title:         video.Title,
				PublishTime:   video.PublishTime.Format("2006-01-02 15:04"),
				Topic:         video.Topic,
			}
			if _, ok := favoriteMap[video.ID]; ok {
				vv.IsFavorite = true
			}
			videoResponse = append(videoResponse, vv)
		} else {
			err := errors.New(constant.VideoServerBug)
			zap.L().Error(err.Error())
			return nil, err
		}
	}
	return videoResponse, nil
}
