package storage

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"diktok/package/constant"
	"diktok/storage/database/model"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var CommentRedis, VideoRedis *redis.Client

// 评论增加 会影响视频的评论数 和评论表 需要lua脚本保证原子性 （目前采取删缓存）
// 评论列表zset吧 按照评论时间排序（可以考虑时间加赞数加权排序）
func CommentAdd(c *model.Comment) error {
	zsetKey := constant.CommentPrefix + strconv.FormatInt(c.VideoID, 10)
	// 应该删缓存 而不是增加 有过期时间的 过期了怎么办
	// 更新倒是可以考虑 但是可能有数据不一致的情况
	err := CommentRedis.Del(zsetKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	videoCountKey := constant.VideoInfoCountPrefix + strconv.FormatInt(c.VideoID, 10)
	err = VideoRedis.Del(videoCountKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

// lua脚本保证原子性 （目前采取删缓存）
func CommentDelete(c *model.Comment) error {
	zsetKey := constant.CommentPrefix + strconv.FormatInt(c.VideoID, 10)
	dataJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = CommentRedis.ZRem(zsetKey, dataJSON).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	// 减少视频的评论数
	videoCountKey := constant.VideoInfoCountPrefix + strconv.FormatInt(c.VideoID, 10)
	err = VideoRedis.Del(videoCountKey).Err()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func SetComments(videoID int64, comments []*model.Comment) error {
	zsetKey := constant.CommentPrefix + strconv.FormatInt(videoID, 10)
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
	pp := CommentRedis.Pipeline()
	pp.ZAdd(zsetKey, members...).Err()
	pp.Expire(zsetKey, constant.Expiration+time.Duration(rand.Intn(200))*time.Second)
	_, err := pp.Exec()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func GetCommentsByVideoID(videoID int64) ([]*model.Comment, error) {
	zsetKey := constant.CommentPrefix + strconv.FormatInt(videoID, 10)
	commentsJSON, err := CommentRedis.ZRevRange(zsetKey, 0, -1).Result()
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
