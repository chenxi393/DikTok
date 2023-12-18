package cache

import (
	"douyin/config"
	"douyin/model"
	"douyin/package/constant"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var VideoRedisClient *redis.Client
var VideoIDBloomFilter *bloom.BloomFilter

// VideoInfo 视频固定的信息
type VideoInfo struct {
	ID          uint64
	PublishTime time.Time
	AuthorID    uint64
	PlayURL     string
	CoverURL    string
	Title       string
}

func InitVideoRedis() {
	// videoRedis 连接
	VideoRedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.System.VideoRedis.Host, config.System.VideoRedis.Port),
		Password: config.System.VideoRedis.Password,
		DB:       config.System.VideoRedis.Database,
		PoolSize: config.System.VideoRedis.PoolSize, //每个CPU最大连接数
	})
	_, err := VideoRedisClient.Ping().Result()
	if err != nil {
		zap.L().Fatal("video_redis连接失败", zap.Error(err))
	}
	initVideoBloomFilter()
}

// 初始化布隆过滤器
// 布隆过滤器的预估元素数量 和误报率 决定了底层bitmap的大小 和 无偏哈希函数的个数
func initVideoBloomFilter() {
	VideoIDBloomFilter = bloom.NewWithEstimates(100000, 0.01)
	videoIDList := make([]uint64, 0)
	constant.DB.Model(&model.Video{}).Select("id").Find(&videoIDList)
	for _, v := range videoIDList {
		VideoIDBloomFilter.AddString(strconv.FormatUint(v, 10))
	}
	zap.L().Info("初始化布隆过滤器: 成功")
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
	pubulishIDStrings := make([]string, 1, len(pubulishIDSet)+1)
	pubulishIDStrings[0] = "0"
	for i := range pubulishIDSet {
		pubulishIDStrings = append(pubulishIDStrings, strconv.FormatUint(pubulishIDSet[i], 10))
	}
	pp := VideoRedisClient.Pipeline()
	pp.SAdd(key, pubulishIDStrings)
	pp.Expire(key, constant.Expiration+time.Duration(rand.Intn(100))*time.Second)
	_, err := pp.Exec()
	return err
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

func PublishVideo(userID, videoID uint64) error {
	publishSet := constant.PublishIDPrefix + strconv.FormatUint(userID, 10)
	authorInfoCountKey := constant.UserInfoCountPrefix + strconv.FormatUint(userID, 10)
	authorInfoKey := constant.UserInfoPrefix + strconv.FormatUint(userID, 10)
	// 这里也应该删缓存 不能增加
	err := UserRedisClient.Del(authorInfoCountKey, authorInfoKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	// 应该删缓存 而不是增加 或者重新设置整个集合
	// TODO可以考虑把视频加入 以便feed使用
	err = VideoRedisClient.Del(publishSet).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}
