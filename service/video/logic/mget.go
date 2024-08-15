package logic

import (
	"context"
	"sync"
	"time"

	"diktok/config"
	pbfavorite "diktok/grpc/favorite"
	pbuser "diktok/grpc/user"
	pbvideo "diktok/grpc/video"
	"diktok/package/rpc"
	stroage "diktok/service/video/storage"
	"diktok/storage/database"
	"diktok/storage/database/model"
	"diktok/storage/database/query"

	"go.uber.org/zap"
	"gorm.io/gen"
)

func MGetVideos(ctx context.Context, req *pbvideo.MGetReq) (*pbvideo.MGetResp, error) {
	var offset int32 = 0
	if req.Offset > 0 {
		offset = req.GetOffset()
	}
	var limit int32 = 50
	if req.Limit > 0 {
		limit = req.GetLimit()
	}
	so := query.Use(database.DB).Video
	var conds []gen.Condition
	if req.UserId > 0 {
		conds = append(conds, so.AuthorID.Eq(req.GetUserId()))
	}
	if len(req.VideoId) > 0 {
		conds = append(conds, so.ID.In(req.GetVideoId()...))
	}
	if req.MaxPublishTime > 0 {
		conds = append(conds, so.PublishTime.Lte(time.UnixMilli(req.GetMaxPublishTime())))
	}
	if req.Topic != "" {
		conds = append(conds, so.Topic.Like(req.Topic+"%"))
	}
	// 默认按ID降序排序
	videos, total, err := stroage.MGetVideosByCond(ctx, int(offset), int(limit), conds, so.ID.Desc())
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &pbvideo.MGetResp{
		VideoList: buildMGetVideosResp(videos),
		Total:     total,
		HasMode:   total > int64(offset+limit),
	}, nil
}

func buildMGetVideosResp(videos []*model.Video) []*pbvideo.VideoMetaData {
	res := make([]*pbvideo.VideoMetaData, 0, len(videos))
	for _, v := range videos {
		res = append(res, &pbvideo.VideoMetaData{
			Id:          v.ID,
			AuthorId:    v.AuthorID,
			PlayUrl:     config.System.Qiniu.OssDomain + "/" + v.PlayURL,
			CoverUrl:    config.System.Qiniu.OssDomain + "/" + v.CoverURL,
			Title:       v.Title,
			Topic:       v.Topic,
			PublishTime: v.PublishTime.UnixMilli(), //Format("2006-01-02 15:04"), // 这个时间戳 应该给前端转换 TODO
		})
	}
	return res
}

// TODO 与 user 点赞 解耦
// RPC调用拿userMap 里的用户信息 拿video里的详细信息 返回
// 然后 对这些视频判断 loginUserID 有没有点赞
// 视频里的作者信息应当都是service.user_id（还需判断 登录用户有没有关注）
func getVideoInfo(ctx context.Context, videos []*model.Video, userMap map[int64]*pbuser.UserInfo, loginUserID int64) []*pbvideo.VideoData {
	// rpc调用 去拿个人信息
	wg := &sync.WaitGroup{}
	wg.Add(len(userMap))
	for userID := range userMap {
		go func(id int64) {
			defer wg.Done()
			// TODO 这里是不是也应该 rpc批量拿出来 而不是一个个去拿
			user, err := rpc.UserClient.Info(ctx, &pbuser.InfoRequest{
				LoginUserID: loginUserID,
				UserID:      id,
			})
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			if user.StatusCode != 0 {
				zap.L().Error("rpc 调用错误")
			}
			*userMap[id] = *user.GetUser()
		}(userID)
	}
	videoInfos := make([]*pbvideo.VideoData, 0, len(videos))
	for i := range videos {
		videoInfos = append(videoInfos, videoDataInfo(videos[i], userMap[videos[i].AuthorID]))
	}
	// 判断请求用户是否喜欢
	if loginUserID != 0 {
		wg.Add(len(videoInfos))
		for i := range videoInfos {
			go func(i int) {
				defer wg.Done()
				result, err := rpc.FavoriteClient.IsFavorite(ctx, &pbfavorite.IsFavoriteRequest{
					UserID:  loginUserID,
					VideoID: videoInfos[i].Id,
				})
				if err != nil {
					zap.L().Error(err.Error())
					return
				}
				videoInfos[i].IsFavorite = result.IsFavorite
			}(i)
		}
	}
	wg.Wait()
	return videoInfos
}

func videoDataInfo(v *model.Video, u *pbuser.UserInfo) *pbvideo.VideoData {
	return &pbvideo.VideoData{
		Author:        u,
		Id:            v.ID,
		PlayUrl:       config.System.Qiniu.OssDomain + "/" + v.PlayURL,
		CoverUrl:      config.System.Qiniu.OssDomain + "/" + v.CoverURL,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		Title:         v.Title,
		Topic:         v.Topic,
		PublishTime:   v.PublishTime.Format("2006-01-02 15:04"),
	}
}
