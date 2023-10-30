package cache

import (
	"douyin/model"
	"douyin/package/constant"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// VideoInfo 视频固定的信息
type VideoInfo struct {
	ID          uint64
	PublishTime time.Time
	AuthorID    uint64
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
	pipeline := VideoRedisClient.Pipeline()

	// 下面两个的过期时间保持一致 不然查库还是会查出信息
	randomTime := constant.Expiration + time.Duration(rand.Intn(100))*time.Second
	// 设置 UserInfo 的 JSON 缓存
	infoKey := constant.VideoInfoPrefix + strconv.FormatUint(video.ID, 10)
	err = pipeline.Set(infoKey, videoInfoJSON, randomTime).Err()
	if err != nil {
		zap.L().Sugar().Error(err)
		return err
	}

	infoCountKey := constant.VideoInfoCountPrefix + strconv.FormatUint(video.ID, 10)
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

func GetVideoInfo(videoID uint64) (*model.Video, error) {
	infoKey := constant.VideoInfoPrefix + strconv.FormatUint(videoID, 10)
	infoCountKey := constant.VideoInfoCountPrefix + strconv.FormatUint(videoID, 10)
	// 使用管道加速
	pipeline := VideoRedisClient.Pipeline()
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

func SetPublishSet(userID uint64, pubulishIDSet []uint64) error {
	key := constant.PublishIDPrefix + strconv.FormatUint(userID, 10)
	pubulishIDStrings := make([]string, 0, len(pubulishIDSet))
	for i := range pubulishIDSet {
		pubulishIDStrings = append(pubulishIDStrings, strconv.FormatUint(pubulishIDSet[i], 10))
	}
	return VideoRedisClient.SAdd(key, pubulishIDStrings).Err()
}

func GetPubulishSet(userID uint64) ([]uint64, error) {
	key := constant.PublishIDPrefix + strconv.FormatUint(userID, 10)
	idSet, err := VideoRedisClient.SMembers(key).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	res := make([]uint64, 0, len(idSet))
	for _, t := range idSet {
		id, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		res = append(res, id)
	}
	return res, nil
}

// func PublishAction(userID, videoID uint64, cnt int64) error {
// 	videoInfoCountKey := constant.VideoInfoCountPrefix + strconv.FormatUint(videoID, 10)
// 	favoriteKey := constant.FavoriteIDPrefix + strconv.FormatUint(userID, 10)
// 	userInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(userID, 10)
// 	authorInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(author_id, 10)
// 	pp := UserRedisClient.Pipeline()
// 	if cnt == 1 {
// 		pp.SAdd(favoriteKey, strconv.FormatUint(videoID, 10))
// 	} else {
// 		pp.SRem(favoriteKey, strconv.FormatUint(videoID, 10))
// 	}
// 	pp.HIncrBy(videoInfoCountKey, constant.FavoritedCountField, cnt)
// 	pp.HIncrBy(userInfoCountKey, constant.FavoriteCountField, cnt)
// 	pp.HIncrBy(authorInfoCountKey, constant.TotalFavoritedField, cnt)
// 	_, err := pp.Exec()
// 	if err != nil {
// 		zap.L().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }