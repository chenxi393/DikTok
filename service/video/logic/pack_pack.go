package logic

import (
	"context"

	"diktok/config"
	pbcomment "diktok/grpc/comment"
	pbfavorite "diktok/grpc/favorite"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/rpc"
	"diktok/package/util"
	"time"

	"github.com/sourcegraph/conc"
	"go.uber.org/zap"
)

// copy 自gate way
// 可以作为通用的 视频id 和登录用户打包 视频信息
// 打包服务 包括 load 和 pack 应该分开
func BuildVideosInfo(ctx context.Context, videoIDs []int64, videoMeta []*pbvideo.VideoMetaData, loginUserID int64) (resp []*pbvideo.VideoData, err error) {
	if len(videoIDs) <= 0 && len(videoMeta) <= 0 {
		return nil, nil
	} else if len(videoIDs) <= 0 && len(videoMeta) > 0 {
		for _, v := range videoMeta {
			videoIDs = append(videoIDs, v.Id)
		}
	}
	var likeMap map[int64]bool
	var likeCount map[int64]int64
	var commentCount map[int64]int64
	var UserMap map[int64]*pbuser.UserInfo
	var wg conc.WaitGroup
	// 是否点赞
	wg.Go(func() {
		// var wgg conc.WaitGroup
		// mu := sync.Mutex{}
		for _, v := range videoIDs {
			// wgg.Go(func() {
			// 这里有没有办法批量判断 TODO 或者拿登录用户的点赞视频列表？ 并发数量很大经常报错
			result, err := rpc.FavoriteClient.IsFavorite(ctx, &pbfavorite.IsFavoriteRequest{
				UserID:  loginUserID,
				VideoID: v,
			})
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			// mu.Lock()
			likeMap[v] = result.GetIsFavorite()
			// mu.Unlock()
			// })
		}
		// wgg.WaitAndRecover()
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
		resp, err := rpc.CommentClient.Count(ctx, &pbcomment.CountReq{
			ParentIDs: videoIDs,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		commentCount = resp.GetCountMap()
	})
	wg.WaitAndRecover()
	userIDs := make([]int64, 0, len(videoMeta))
	for _, v := range videoMeta {
		userIDs = append(userIDs, v.AuthorId)
	}
	userResp, err := rpc.UserClient.List(ctx, &pbuser.ListReq{
		UserID: userIDs,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	UserMap = userResp.GetUser()
	zap.L().Sugar().Infof("[BuildVideosInfo] videoMeta = %s,UserMap = %s", util.GetLogStr(videoMeta), util.GetLogStr(UserMap))
	return buildVideoInfo(videoMeta, UserMap, likeMap, likeCount, commentCount), nil
}

func buildVideoInfo(items []*pbvideo.VideoMetaData, userMap map[int64]*pbuser.UserInfo, isLiked map[int64]bool, likeCount map[int64]int64, commentCount map[int64]int64) []*pbvideo.VideoData {
	data := make([]*pbvideo.VideoData, 0, len(items))
	for _, item := range items {
		v := &pbvideo.VideoData{
			Author:        buildUserInfo(userMap[item.AuthorId]),
			Id:            item.Id,
			PlayUrl:       config.System.Qiniu.OssDomain + "/" + item.PlayUrl,
			CoverUrl:      config.System.Qiniu.OssDomain + "/" + item.CoverUrl,
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

func buildUserInfo(user *pbuser.UserInfo) *pbuser.UserInfo {
	if user == nil {
		return nil
	}
	return &pbuser.UserInfo{
		Avatar:          config.System.Qiniu.OssDomain + "/" + user.Avatar,
		BackgroundImage: config.System.Qiniu.OssDomain + "/" + user.BackgroundImage,
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Id:              user.Id,
		IsFollow:        user.IsFollow, // TODO 是否关注这里还是交给下游聚合把？？
		Name:            user.Name,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
}
