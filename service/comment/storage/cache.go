package storage

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"diktok/package/constant"
	"diktok/package/util"
	"diktok/storage/database/model"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// TODO 找个缓存的框架使用
// TODO 订阅binlog 删缓存
var CommentRedis *redis.Client

// 评论内容缓存
func GetContentByCache(ctx context.Context, commentIDs []int64) ([]*model.CommentContent, error) {
	// redis 查数据
	commentIDsKey := make([]string, 0, len(commentIDs))
	for _, v := range commentIDs {
		commentIDsKey = append(commentIDsKey, fmt.Sprintf(constant.CommentContentPrefix, v))
	}
	commentsCache, err := CommentRedis.MGet(commentIDsKey...).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	// 回源
	sourceIDs := make([]int64, 0)
	for i, v := range commentsCache {
		if v == nil {
			sourceIDs = append(sourceIDs, commentIDs[i])
		}
	}
	var commentsSourceMp map[int64]*model.CommentContent
	if len(sourceIDs) > 0 {
		commentsSource, err := GetContentByIDs(ctx, sourceIDs)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		commentsSourceMp = make(map[int64]*model.CommentContent, len(commentsSource))
		for _, v := range commentsSource {
			commentsSourceMp[v.ID] = v
		}
	}

	// 聚合
	comments := make([]*model.CommentContent, 0, len(commentIDs))
	for i, v := range commentsCache {
		var data *model.CommentContent
		if v == nil {
			data = commentsSourceMp[commentIDs[i]]
		} else {
			zap.L().Sugar().Debugf("[GetContentByCache] data = %s", util.GetLogStr(v))
			data = &model.CommentContent{}
			err = sonic.Unmarshal([]byte(v.(string)), data)
			if err != nil {
				return nil, err
			}
		}
		comments = append(comments, data)
	}
	// 异步写缓存
	// 注意这里无论异步还是同步 都存在缓存不一致的问题
	// 例如上面 拿了数据 这时 数据被更改 MQ消费删缓存消息 在这次异步写入之前
	go func() {
		// 使用管道批量执行命令
		pipe := CommentRedis.Pipeline()

		for _, v := range commentsSourceMp {
			contentJson, _ := sonic.Marshal(v)
			pipe.Set(fmt.Sprintf(constant.CommentContentPrefix, v.ID), contentJson, constant.Expiration+time.Duration(rand.Intn(50))*time.Second)
		}

		// 执行管道中的命令
		_, err := pipe.Exec()
		if err != nil {
			zap.L().Error("Error executing pipeline: " + err.Error())
			return
		}
	}()
	zap.L().Sugar().Debugf("[GetContentByCache] commentIDs = %s, sourceIDs = %s, comments = %s", util.GetLogStr(commentIDs), util.GetLogStr(sourceIDs), util.GetLogStr(comments))
	return comments, nil
}

// // lua脚本保证原子性 （目前采取删缓存）
// func CommentDelete(c *model.Comment) error {
// 	zsetKey := constant.CommentPrefix + strconv.FormatInt(c.VideoID, 10)
// 	dataJSON, err := json.Marshal(c)
// 	if err != nil {
// 		return err
// 	}
// 	err = CommentRedis.ZRem(zsetKey, dataJSON).Err()
// 	if err != nil {
// 		zap.L().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// TODO 获取一个视频 评论元信息
// func GetCommentsByVideoID(videoID int64) ([]*model.Comment, error) {
// 	zsetKey := constant.CommentPrefix + strconv.FormatInt(videoID, 10)
// 	commentsJSON, err := CommentRedis.ZRevRange(zsetKey, 0, -1).Result()
// 	if err != nil {
// 		zap.L().Error(err.Error())
// 		return nil, err
// 	}
// 	// ZRange 查不到数据不会返回 redis.Nil
// 	if len(commentsJSON) == 0 {
// 		return nil, redis.Nil
// 	}
// 	comments := make([]*model.Comment, 0, len(commentsJSON))
// 	for _, id := range commentsJSON {
// 		var data model.Comment
// 		err = json.Unmarshal([]byte(id), &data)
// 		if err != nil {
// 			return nil, err
// 		}
// 		comments = append(comments, &data)
// 	}
// 	return comments, nil
// }

// func SetComments(videoID int64, comments []*model.Comment) error {
// 	zsetKey := constant.CommentPrefix + strconv.FormatInt(videoID, 10)
// 	members := make([]redis.Z, 0, len(comments))
// 	for _, c := range comments {
// 		dataJSON, err := json.Marshal(c)
// 		if err != nil {
// 			zap.L().Error(err.Error())
// 			return err
// 		}
// 		member := redis.Z{
// 			Score:  float64(c.CreatedTime.UnixMilli()),
// 			Member: dataJSON,
// 		}
// 		members = append(members, member)
// 	}
// 	pp := CommentRedis.Pipeline()
// 	pp.ZAdd(zsetKey, members...).Err()
// 	pp.Expire(zsetKey, constant.Expiration+time.Duration(rand.Intn(200))*time.Second)
// 	_, err := pp.Exec()
// 	if err != nil {
// 		zap.L().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }
