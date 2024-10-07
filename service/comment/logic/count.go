package logic

import (
	"context"

	pbcomment "diktok/grpc/comment"
	"diktok/service/comment/storage"
	"diktok/storage/database"
	"diktok/storage/database/query"

	"go.uber.org/zap"
	"gorm.io/gen"
)

func Count(ctx context.Context, req *pbcomment.CountReq) (*pbcomment.CountResp, error) {
	so := query.Use(database.DB).CommentMetum
	var conds []gen.Condition
	if req.ItemIdIndex == 0 {
		conds = append(conds, so.ItemID.In(req.ParentIDs...))
	} else {
		conds = append(conds, so.ItemID.Eq(req.ItemIdIndex))
		conds = append(conds, so.ParentID.In(req.ParentIDs...)) // 查询子评论
	}
	countMap, err := storage.CountMapByIDs(ctx, conds)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &pbcomment.CountResp{
		CountMap: countMap,
	}, nil
}
