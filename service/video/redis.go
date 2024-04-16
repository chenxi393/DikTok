package main

import (
	"douyin/package/constant"
	"douyin/storage/database/model"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// VideoInfo 视频固定的信息
type VideoInfo struct {
	ID          int64
	PublishTime time.Time
	AuthorID    int64
	PlayURL     string
	CoverURL    string
	Title       string
}

func SetVideoInfo(video *model.Video) error {
	videoInfo := &VideoInfo{
		ID:          video.ID,
		PublishTime: video.PublishTime,
		AuthorID:    video.AuthorID,
		PlayURL:     video.PlayURL,
		CoverURL:    video.CoverURL,
		Title:       video.Title,
	}
	videoInfoJSON, err := json.Marshal(videoInfo)
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	// 开启管道 一次发送请求
	pipeline := videoRedis.Pipeline()

	// 下面两个的过期时间保持一致 不然查库还是会查出信息
	randomTime := constant.Expiration + time.Duration(rand.Intn(100))*time.Second
	// 设置 UserInfo 的 JSON 缓存
	infoKey := constant.VideoInfoPrefix + strconv.FormatInt(video.ID, 10)
	err = pipeline.Set(infoKey, videoInfoJSON, randomTime).Err()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}

	infoCountKey := constant.VideoInfoCountPrefix + strconv.FormatInt(video.ID, 10)
	// 使用 MSet 进行批量设置
	err = pipeline.HMSet(infoCountKey, map[string]interface{}{
		constant.FavoritedCountField: video.FavoriteCount,
		constant.CommentCountField:   video.CommentCount,
	}).Err()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	err = pipeline.Expire(infoCountKey, randomTime).Err()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	// 执行管道中的命令
	_, err = pipeline.Exec()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}
	return nil
}

func GetVideoInfo(videoID int64) (*model.Video, error) {
	infoKey := constant.VideoInfoPrefix + strconv.FormatInt(videoID, 10)
	infoCountKey := constant.VideoInfoCountPrefix + strconv.FormatInt(videoID, 10)
	// 使用管道加速
	pipeline := videoRedis.Pipeline()
	// 注意pipeline返回指针 返回值肯定是nil
	videoInfoCmd := pipeline.Get(infoKey)
	videoInfoCountCmd := pipeline.HGetAll(infoCountKey)
	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}
	// 提取返回的结果
	videoInfo, err := videoInfoCmd.Result()
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}
	videoInfoCount, err := videoInfoCountCmd.Result()
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}
	// 解析不变的字段
	videoInfoFixed := VideoInfo{}
	err = json.Unmarshal([]byte(videoInfo), &videoInfoFixed)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}

	// 解析count信息
	favoriteCount, _ := strconv.ParseInt(videoInfoCount[constant.FavoritedCountField], 10, 64)
	commentCount, _ := strconv.ParseInt(videoInfoCount[constant.CommentCountField], 10, 64)

	return &model.Video{
		ID:            videoInfoFixed.ID,
		AuthorID:      videoInfoFixed.AuthorID,
		PlayURL:       videoInfoFixed.PlayURL,
		CoverURL:      videoInfoFixed.CoverURL,
		Title:         videoInfoFixed.Title,
		PublishTime:   videoInfoFixed.PublishTime,
		FavoriteCount: favoriteCount,
		CommentCount:  commentCount,
	}, nil
}

func SetPublishSet(userID int64, pubulishIDSet []int64) error {
	key := constant.PublishIDPrefix + strconv.FormatInt(userID, 10)
	pubulishIDStrings := make([]string, 1, len(pubulishIDSet)+1)
	pubulishIDStrings[0] = "0"
	for i := range pubulishIDSet {
		pubulishIDStrings = append(pubulishIDStrings, strconv.FormatInt(pubulishIDSet[i], 10))
	}
	pp := videoRedis.Pipeline()
	pp.SAdd(key, pubulishIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
}

func GetPubulishSet(userID int64) ([]int64, error) {
	key := constant.PublishIDPrefix + strconv.FormatInt(userID, 10)
	idSet, err := videoRedis.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	res := make([]int64, 0, len(idSet))
	for _, t := range idSet {
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

func PublishVideo(userID, videoID int64) error {
	publishSet := constant.PublishIDPrefix + strconv.FormatInt(userID, 10)
	authorInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatInt(userID, 10)
	authorInfoKey := constant.UserInfoPrefix + strconv.FormatInt(userID, 10)
	// 这里也应该删缓存 不能增加
	err := userRedis.Del(authorInfoCountKey, authorInfoKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	// 应该删缓存 而不是增加 或者重新设置整个集合
	// TODO可以考虑把视频加入 以便feed使用
	err = videoRedis.Del(publishSet).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
