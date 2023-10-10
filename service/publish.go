package service

import (
	"bytes"
	"douyin/config"
	"douyin/dal/dao"
	"douyin/dal/model"
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
	// 2: optional binary data // 视频数据
	Title string `form:"title"`
}

func (service *PublisService) PublishAction(userID uint64, buf []byte) (*response.PublishResponse, error) {
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
	_, err = dao.CreateVideo(&model.Video{
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
	return &response.PublishResponse{
		StatusCode: response.Success,
		StatusMsg:  "上传视频成功",
	}, nil
}
