package logic

import (
	"context"
	"time"

	"diktok/config"
	pbvideo "diktok/grpc/video"
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
