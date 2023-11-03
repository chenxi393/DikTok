package service

import (
	"bytes"
	"douyin/config"
	"douyin/database"
	"douyin/model"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/util"
	"douyin/response"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/plugin/dbresolver"
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

func (service *PublisService) PublishAction(userID uint64, buf *bytes.Buffer) (*response.CommonResponse, error) {
	u1, err := uuid.NewV4()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	fileName := u1.String() + "." + "mp4"
	// 这里暂时上传到本地
	playURL, coverURL, err := util.UploadVideo(buf.Bytes(), fileName)
	if err != nil {
		return nil, err
	}
	video_id, err := database.CreateVideo(&model.Video{
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
	//加入布隆过滤器
	cache.VideoIDBloomFilter.AddString(strconv.FormatUint(video_id, 10))
	// 异步上传到对象存储
	go func() {
		playURL, err := util.UploadToOSS(fileName, config.System.HttpAddress.VideoAddress+"/"+fileName)
		if err != nil {
			// 这里要不要重试？？
			zap.L().Error(err.Error())
		}
		if coverURL != config.System.HttpAddress.DefaultCoverURL {
			coverFileName := fileName + ".jpeg"
			coverURL, err = util.UploadToOSS(coverFileName, config.System.HttpAddress.ImageAddress+"/"+coverFileName)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}
		// 先更新数据库 再更新缓存
		err = database.UpdateVideoURL(playURL, coverURL, video_id)
		if err != nil {
			zap.L().Error(err.Error())
		}
		// 这里会有主从复制延时导致缓存不一致的问题。。
		// 对于即时写即时读的要指定主库去读 不能读从库
		var video model.Video
		err = constant.DB.Clauses(dbresolver.Write).Where("id = ?", video_id).First(&video).Error
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		cache.SetVideoInfo(&video)
		// TODO 记得删除本地的视频和图片
	}()
	return &response.CommonResponse{
		StatusCode: response.Success,
		StatusMsg:  response.UploadVideoSuccess,
	}, nil
}

func (service *PublishListService) GetPublishVideos(loginUserID uint64) (*response.VideoListResponse, error) {
	// 第一步查找 所有的 service.user_id 的视频记录
	// 然后 对这些视频判断 loginUserID 有没有点赞
	// 视频里的作者信息应当都是service.user_id（还需判断 登录用户有没有关注）
	// TODO 加分布式锁 redis
	// TODO 这里其实应当先去redis拿列表 再去数据库拿数据
	videos, err := database.SelectVideosByUserID(service.UserID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 不都是一个作者嘛 拿一次信息不就好了
	author, err := cache.GetUserInfo(service.UserID)
	if err != nil {
		zap.L().Warn(constant.CacheMiss, zap.Error(err))
		author, err = database.SelectUserByID(service.UserID)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		go func() {
			err := cache.SetUserInfo(author)
			if err != nil {
				zap.L().Error(err.Error())
			}
		}()
	}
	var isFollowed bool
	if loginUserID == 0 {
		isFollowed = false
	}
	if service.UserID == loginUserID {
		isFollowed = true
	} else {
		isFollowed, err = cache.IsFollow(loginUserID, service.UserID)
		if err != nil {
			zap.L().Warn(constant.CacheMiss)
			isFollowed, err = database.IsFollowed(loginUserID, service.UserID)
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
			go func() {
				following, err := database.SelectFollowingByUserID(loginUserID)
				if err != nil {
					zap.L().Error(err.Error())
					return
				}
				err = cache.SetFollowUserIDSet(loginUserID, following)
				if err != nil {
					zap.L().Error(err.Error())
				}
			}()
		}
	}
	var favorite []uint64
	if loginUserID != 0 {
		favorite, err = cache.GetFavoriteSet(loginUserID)
		if err != nil {
			zap.L().Warn(constant.CacheMiss)
			favorite, err = database.SelectFavoriteVideoByUserID(loginUserID)
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
			go func() {
				err := cache.SetFavoriteSet(loginUserID, favorite)
				if err != nil {
					zap.L().Error(err.Error())
				}
			}()
		}
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
		if _, ok := favoriteMap[ff.ID]; ok {
			item.IsFavorite = true
		}
		reps = append(reps, item)
	}
	return &response.VideoListResponse{
		StatusCode: response.Success,
		StatusMsg:  response.PubulishListSuccess,
		VideoList:  reps,
	}, nil
}
