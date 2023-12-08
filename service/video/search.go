package service

import (
	"douyin/config"
	"douyin/database"
	"douyin/model"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/response"
	"errors"

	"go.uber.org/zap"
)

type SearchService struct {
	// 搜索框输入
	KeyWord string `query:"keyword"`
	// 用户登录状态下设置
	Token string `query:"token"`
}

func (service *SearchService) SearchVideo(userID uint64) (*response.VideoListResponse, error) {
	if service.KeyWord == "" {
		return nil, errors.New(constant.BadParaRequest)
	}
	// 去数据库利用全文索引拿出所有视频数据
	videos, err := database.SearchVideoByKeyword(service.KeyWord)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	// 拿到视频数据之后 还得一个视频一个视频拿到作者信息
	userIDs := make([]uint64, 0, len(videos))
	for _, video := range videos {
		userIDs = append(userIDs, video.AuthorID)
	}
	// TODO 作者信息应该先去redis里面拿 没有再数据库填空
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
			_, isFollow := followingMap[video.AuthorID]
			vv := response.Video{
				Author:        *response.UserInfo(usersMap[video.AuthorID], isFollow),
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
	return &response.VideoListResponse{
		StatusCode: response.Success,
		StatusMsg:  response.ActionSuccess,
		VideoList:  videoResponse,
	}, nil

}
