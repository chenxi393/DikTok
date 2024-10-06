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
	countMap := make(map[int64]int64, len(req.GetParentIDs()))
	so := query.Use(database.DB).CommentMetum
	// FIXME 循环SQL 遭不住 这里是不是可以优化SQL 一次给查出来
	for _, v := range req.GetParentIDs() {
		var conds []gen.Condition
		if req.ItemIdIndex == 0 {
			conds = append(conds, so.ItemID.Eq(v))
		} else {
			conds = append(conds, so.ItemID.Eq(req.ItemIdIndex))
			conds = append(conds, so.ParentID.Eq(v)) // 查询子评论
		}
		total, err := storage.CountByCond(ctx, conds)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		countMap[v] = total
	}
	return &pbcomment.CountResp{
		CountMap: countMap,
	}, nil
}
