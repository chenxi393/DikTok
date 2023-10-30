package cache

import (
	"douyin/model"
	"douyin/package/constant"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// FIXME评论增加 会影响视频的评论数 和评论表 需要lua脚本保证原子性
// 评论列表zset吧 按照评论时间排序（可以考虑时间加赞数加权排序）
func CommentAdd(c *model.Comment) error {
	zsetKey := constant.CommentPrefix + strconv.FormatUint(c.VideoID, 10)
	dataJSON, err := json.Marshal(c)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	member := redis.Z{
		Score:  float64(c.CreatedTime.UnixMilli()),
		Member: dataJSON,
	}
	pp := CommentRedisClient.Pipeline()
	pp.ZAdd(zsetKey, member).Err()
	pp.Expire(zsetKey, constant.Expiration+time.Duration(rand.Intn(200))*time.Second)
	// 增加视频的评论数
	videoCountKey := constant.VideoInfoCountPrefix + strconv.FormatUint(c.VideoID, 10)
	err = VideoRedisClient.HIncrBy(videoCountKey, constant.CommentCountField, 1).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	_, err = pp.Exec()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

// FIXME lua脚本保证原子性
func CommentDelete(c *model.Comment) error {
	zsetKey := constant.CommentPrefix + strconv.FormatUint(c.VideoID, 10)
	dataJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = CommentRedisClient.ZRem(zsetKey, dataJSON).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	// 减少视频的评论数
	videoCountKey := constant.VideoInfoCountPrefix + strconv.FormatUint(c.VideoID, 10)
	err = VideoRedisClient.HIncrBy(videoCountKey, constant.CommentCountField, -1).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func SetComments(videoID uint64, comments []*model.Comment) error {
	zsetKey := constant.CommentPrefix + strconv.FormatUint(videoID, 10)
	members := make([]redis.Z, 0, len(comments))
	for _, c := range comments {
		dataJSON, err := json.Marshal(c)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		member := redis.Z{
			Score:  float64(c.CreatedTime.UnixMilli()),
			Member: dataJSON,
		}
		members = append(members, member)
	}
	pp := CommentRedisClient.Pipeline()
	pp.ZAdd(zsetKey, members...).Err()
	pp.Expire(zsetKey, constant.Expiration+time.Duration(rand.Intn(200))*time.Second)
	_, err := pp.Exec()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func GetCommentsByVideoID(videoID uint64) ([]*model.Comment, error) {
	zsetKey := constant.CommentPrefix + strconv.FormatUint(videoID, 10)
	commentsJSON, err := CommentRedisClient.ZRange(zsetKey, 0, -1).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// ZRange 查不到数据不会返回 redis.Nil
	if len(commentsJSON) == 0 {
		return nil, redis.Nil
	}
	comments := make([]*model.Comment, 0, len(commentsJSON))
	for _, id := range commentsJSON {
		var data model.Comment
		err = json.Unmarshal([]byte(id), &data)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &data)
	}
	return comments, nil
}
