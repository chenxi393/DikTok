package logic

import (
	"context"
	"errors"

	pbcomment "diktok/grpc/comment"
	pbfavorite "diktok/grpc/favorite"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"
	"time"

	"github.com/sourcegraph/conc"
	"go.uber.org/zap"
)

// 作为通用的 视频id（视频元信息） 和登录用户打包视频信息
func Pack(ctx context.Context, req *pbvideo.PackReq) (*pbvideo.PackResp, error) {
	if len(req.GetVideoId()) <= 0 {
		return nil, errors.New(constant.BadParaRequest)
	}

	videoInfo, err := BuildVideosInfo(ctx, req.GetVideoId(), nil, req.GetLoginUserId())
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &pbvideo.PackResp{
		VideoList: videoInfo,
	}, nil
}

// load 和 pack 应该分开
func BuildVideosInfo(ctx context.Context, videoIDs []int64, videoMeta []*pbvideo.VideoMetaData, loginUserID int64) (resp []*pbvideo.VideoData, err error) {
	if len(videoIDs) <= 0 && len(videoMeta) <= 0 {
		return nil, nil
	} else if len(videoIDs) <= 0 && len(videoMeta) > 0 {
		for _, v := range videoMeta {
			videoIDs = append(videoIDs, v.Id)
		}
	} else if len(videoMeta) <= 0 {
		videos, err := rpc.VideoClient.MGet(ctx, &pbvideo.MGetReq{
			VideoId: videoIDs,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		videoMeta = videos.GetVideoList()
		// Mget接口不会按照入参顺序返回 需要手动聚合
		// for _, v := range videoMeta {
		// 	videoMetaDataMap[v.Id] = v
		// }
	}

	// var videoMetaDataMap = make(map[int64]*pbvideo.VideoMetaData)
	var likeMap = make(map[int64]bool)
	var likeCount map[int64]int64
	var commentCount map[int64]int64
	var UserMap map[int64]*pbuser.UserInfo
	var wg conc.WaitGroup
	// 用户信息
	wg.Go(func() {
		userIDs := make([]int64, 0, len(videoMeta))
		for _, v := range videoMeta {
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
	})
	// 是否点赞
	wg.Go(func() {
		if loginUserID == 0 {
			return
		}
		result, err := rpc.FavoriteClient.IsFavorite(ctx, &pbfavorite.IsFavoriteReq{
			UserID:  loginUserID,
			VideoID: videoIDs,
		})
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		likeMap = result.GetIsFavorite()
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
	zap.L().Sugar().Infof("[BuildVideosInfo] videoMeta = %s,UserMap = %s", util.GetLogStr(videoMeta), util.GetLogStr(UserMap))
	return buildVideoInfo(videoMeta, UserMap, likeMap, likeCount, commentCount), nil
}

func buildVideoInfo(items []*pbvideo.VideoMetaData, userMap map[int64]*pbuser.UserInfo, isLiked map[int64]bool, likeCount map[int64]int64, commentCount map[int64]int64) []*pbvideo.VideoData {
	data := make([]*pbvideo.VideoData, 0, len(items))
	for _, item := range items {
		v := &pbvideo.VideoData{
			Author:        userMap[item.AuthorId],
			Id:            item.Id,
			PlayUrl:       util.Uri2Url(item.PlayUri),
			CoverUrl:      util.Uri2Url(item.CoverUri),
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
