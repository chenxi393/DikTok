package logic

import (
	"context"
	"os"
	"strconv"
	"time"

	"diktok/config"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/service/video/storage"
	"diktok/storage/cache"
	"diktok/storage/database"
	"diktok/storage/database/model"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/plugin/dbresolver"
)

func Publish(ctx context.Context, req *pbvideo.PublishRequest) (*pbvideo.PublishResponse, error) {
	// 生成唯一文件名
	u1, err := uuid.NewV4()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	fileName := u1.String() + "." + "mp4"

	playURL, coverURL, err := storage.UploadVideo(req.Data, fileName)
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
	video_id, err := storage.CreateVideo(&model.Video{
		PublishTime:   time.Now(),
		AuthorID:      req.LoginUserId,
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
	//加入布隆过滤器  TODO 移出布隆过滤器
	cache.VideoIDBloomFilter.AddString(strconv.FormatInt(video_id, 10))
	// 异步上传到对象存储 可以改用消息队列
	go func() {
		localVideoPath := config.System.HTTP.VideoAddress + "/" + fileName
		err := storage.UploadToOSS(fileName, localVideoPath)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		coverURL = u1.String() + "." + "jpg"
		err = storage.UpdateVideoURL(playURL, coverURL, video_id)
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
		storage.SetVideoInfo(&video)
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
