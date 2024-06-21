package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"diktok/config"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/storage/cache"
	"diktok/storage/database"
	"diktok/storage/database/model"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/plugin/dbresolver"
)

// userID =0 表示未登录
func (s *VideoService) Publish(ctx context.Context, req *pbvideo.PublishRequest) (*pbvideo.PublishResponse, error) {
	// 生成唯一文件名
	u1, err := uuid.NewV4()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	fileName := u1.String() + "." + "mp4"

	playURL, coverURL, err := uploadVideo(req.Data, fileName)
	if err != nil {
		return nil, err
	}
	switch req.Topic {
	case constant.TopicSport:
	case constant.TopicGame:
	case constant.TopicMusic:
	default:
		req.Topic = constant.TopicDefualt + req.Topic
	}
	video_id, err := CreateVideo(&model.Video{
		PublishTime:   time.Now(),
		AuthorID:      req.UserID,
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         req.Title,
		Topic:         req.Topic,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	//加入布隆过滤器
	cache.VideoIDBloomFilter.AddString(strconv.FormatInt(video_id, 10))
	// 异步上传到对象存储
	go func() {
		localVideoPath := config.System.HTTP.VideoAddress + "/" + fileName
		err := uploadToOSS(fileName, localVideoPath)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		coverURL = u1.String() + "." + "jpg"
		err = UpdateVideoURL(playURL, coverURL, video_id)
		if err != nil {
			zap.L().Error(err.Error())
		}
		// 这里会有主从复制延时导致缓存不一致的问题。。
		// 对于即时写即时读的要指定主库去读 不能读从库
		var video model.Video
		err = database.DB.Clauses(dbresolver.Write).Where("id = ?", video_id).First(&video).Error
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		SetVideoInfo(&video)
		// 删除本地的视频
		err = os.Remove(localVideoPath)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}()
	return &pbvideo.PublishResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.UploadVideoSuccess,
	}, nil
}

func (s *VideoService) List(ctx context.Context, req *pbvideo.ListRequest) (*pbvideo.VideoListResponse, error) {
	// 第一步查找 所有的 service.user_id 的视频记录
	// 然后 对这些视频判断 loginUserID 有没有点赞
	// 视频里的作者信息应当都是service.user_id（还需判断 登录用户有没有关注）
	// TODO 加分布式锁 redis
	// TODO 这里其实应当先去redis拿列表 再去数据库拿数据D
	videos, err := SelectVideosByUserID(req.UserID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 作者都是一个 rpc拿作者信息 作者信息包括关注信息
	userMap := make(map[int64]*pbuser.UserInfo)
	userMap[req.UserID] = &pbuser.UserInfo{}
	videoInfos := getVideoInfo(ctx, videos, userMap, req.LoginUserID)
	return &pbvideo.VideoListResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.PubulishListSuccess,
		VideoList:  videoInfos,
	}, nil
}
