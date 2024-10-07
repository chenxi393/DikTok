package handler

import (
	"context"

	"diktok/config"
	"diktok/gateway/response"
	pbfavorite "diktok/grpc/favorite"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/rpc"
	"sync"
	"time"

	"github.com/sourcegraph/conc"
	"go.uber.org/zap"
)

// 可以作为通用的 视频id 和登录用户打包 视频信息
// 这块逻辑 干掉 合入pack 里 是视频打包
// 应该是 提供视频ID 和login id 然后返回视频详情 统一出口
func BuildVideosInfo(ctx context.Context, videoIDs []int64, videoMeta []*pbvideo.VideoMetaData, loginUserID int64) (resp []*response.Video, err error) {
	if len(videoIDs) <= 0 && len(videoMeta) <= 0 {
		return nil, nil
	} else if len(videoIDs) <= 0 && len(videoMeta) > 0 {
		for _, v := range videoMeta {
			videoIDs = append(videoIDs, v.Id)
		}
	}

	var videoMetaDataMap = make(map[int64]*pbvideo.VideoMetaData)
	var likeMap = make(map[int64]bool)
	var likeCount map[int64]int64
	var commentCount map[int64]int64
	var UserMap map[int64]*pbuser.UserInfo
	var wg conc.WaitGroup
	// 视频源信息
	wg.Go(func() {
		if len(videoMeta) <= 0 {
			videos, err := rpc.VideoClient.MGet(ctx, &pbvideo.MGetReq{
				VideoId: videoIDs,
			})
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			videoMeta = videos.GetVideoList()
		}
		// Mget接口不会按照入参顺序返回 需要手动聚合
		for _, v := range videoMeta {
			videoMetaDataMap[v.Id] = v
		}
	})

	// 是否点赞
	wg.Go(func() {
		var wg conc.WaitGroup
		mu := sync.Mutex{}
		for _, v := range videoIDs {
			wg.Go(func() {
				// 这里有没有办法批量判断 TODO 或者拿登录用户的点赞视频列表？
				result, err := rpc.FavoriteClient.IsFavorite(ctx, &pbfavorite.IsFavoriteRequest{
					UserID:  loginUserID,
					VideoID: v,
				})
				if err != nil {
					zap.L().Error(err.Error())
					return
				}
				mu.Lock()
				likeMap[v] = result.GetIsFavorite()
				mu.Unlock()
			})
		}
		wg.WaitAndRecover()
	})
	// 被赞数量
	wg.Go(func() {
		resp, err := rpc.FavoriteClient.Count(ctx, &pbfavorite.CountReq{
			VideoID: videoIDs,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		likeCount = resp.GetTotal()
	})
	// 评论数量
	wg.Go(func() {
		resp, err := rpc.FavoriteClient.Count(ctx, &pbfavorite.CountReq{
			VideoID: videoIDs,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		commentCount = resp.GetTotal()
	})
	wg.WaitAndRecover()

	userIDs := make([]int64, 0, len(videoMetaDataMap))
	for _, v := range videoMetaDataMap {
		userIDs = append(userIDs, v.AuthorId)
	}
	userResp, err := rpc.UserClient.List(ctx, &pbuser.ListReq{
		UserID:      userIDs,
		LoginUserID: loginUserID,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	UserMap = userResp.GetUser()
	return buildVideoInfo(videoMetaDataMap, UserMap, likeMap, likeCount, commentCount), nil
}

func buildVideoInfo(items map[int64]*pbvideo.VideoMetaData, userMap map[int64]*pbuser.UserInfo, isLiked map[int64]bool, likeCount map[int64]int64, commentCount map[int64]int64) []*response.Video {
	data := make([]*response.Video, 0, len(items))
	for _, item := range items {
		v := &response.Video{
			Author:        buildUserInfo(userMap[item.AuthorId]),
			ID:            item.Id,
			PlayURL:       config.System.Qiniu.OssDomain + "/" + item.PlayUrl,
			CoverURL:      config.System.Qiniu.OssDomain + "/" + item.CoverUrl,
			IsFavorite:    isLiked[item.Id],
			FavoriteCount: likeCount[item.Id],
			CommentCount:  commentCount[item.Id], // 其实按道理 评论和点赞 是视频的属性 但是我们已经拆分了 还是要聚合
			Title:         item.Title,
			PublishTime:   time.UnixMilli(item.PublishTime).Format("2006-01-02 15:04"),
			Topic:         item.Topic,
		}
		data = append(data, v)
	}
	return data
}

func buildUserInfo(user *pbuser.UserInfo) *response.User {
	if user == nil {
		return nil
	}
	return &response.User{
		Avatar:          config.System.Qiniu.OssDomain + "/" + user.Avatar,
		BackgroundImage: config.System.Qiniu.OssDomain + "/" + user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		ID:              user.Id,
		IsFollow:        user.IsFollow, // TODO 是否关注这里还是交给下游聚合把？？
		Name:            user.Name,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
}
