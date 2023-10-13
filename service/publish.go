package service

import (
	"bytes"
	"douyin/config"
	"douyin/database"
	"douyin/model"
	"douyin/package/util"
	"douyin/response"
	"io"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type PublisService struct {
	// 用户鉴权token
	Token string `form:"token"`
	// 视频标题
	Title string `form:"title"`
}

type PublishListService struct {
	// 用户鉴权token
	Token string `query:"token"`
	// 用户id
	UserID uint64 `query:"user_id"`
}

func (service *PublisService) PublishAction(userID uint64, buf []byte) (*response.CommonResponse, error) {
	var reader io.Reader = bytes.NewReader(buf)
	cnt := len(buf)
	u1, err := uuid.NewV4()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	fileName := u1.String() + "." + "mp4"
	var playURL, coverURL string
	if config.SystemConfig.UploadModel == "oss" {
		playURL, coverURL, err = util.UploadVideoToOSS(&reader, cnt, fileName)
	} else {
		playURL, coverURL, err = util.UploadVideoToLocal(&reader, fileName)
	}

	if err != nil {
		return nil, err
	}
	//TODO : 返回被忽略了 需要加入布隆过滤器
	_, err = database.CreateVideo(&model.Video{
		PublishTime:   time.Now(),
		AuthorID:      userID,
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         service.Title,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  "上传视频成功",
	}, nil
}

func (service *PublishListService) GetPublishVideos(loginUserID uint64) (*response.PublishListResponse, error) {
	// TODO 加分布式锁 redis
	// 第一步查找 所有的 service.user_id 的视频记录
	// 然后 对这些视频判断 loginUserID 有没有点赞
	//视频里的作者信息应当都是service.user_id（还需判断 登录用户有没有关注）
	videos, err := database.SelectVideosByUserID(service.UserID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 不都是一个作者嘛 拿一次信息不就好了
	author, err := database.SelectUserByID(service.UserID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	var isFollowed bool
	if service.UserID == loginUserID {
		isFollowed = true
	} else {
		isFollowed, err = database.IsFollowed(loginUserID, service.UserID)
	}
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	favorite, err := database.SelectFavoriteVideoByUserID(service.UserID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	favoriteMap := make(map[uint64]struct{}, len(favorite))
	for _, ff := range favorite {
		favoriteMap[ff] = struct{}{}
	}
	// 构造返回参数
	reps := make([]response.Video, 0, len(videos))
	for i, ff := range videos {
		item := response.Video{
			ID:            videos[i].ID,
			CommentCount:  videos[i].CommentCount,
			CoverURL:      videos[i].CoverURL,
			FavoriteCount: videos[i].FavoriteCount,
			PlayURL:       videos[i].PlayURL,
			Title:         videos[i].Title,
			Author:        *response.UserInfo(author, isFollowed),
		}
		response.UserInfo(author, isFollowed)
		if _, ok := favoriteMap[ff.ID]; ok {
			item.IsFavorite = true
		}
		reps = append(reps, item)
	}
	return &response.PublishListResponse{
		StatusCode: response.Success,
		StatusMsg:  "视频列表加载成功",
		VideoList:  reps,
	}, nil
}
