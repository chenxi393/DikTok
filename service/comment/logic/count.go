package logic

import (
	"context"

	pbcomment "diktok/grpc/comment"
	"diktok/service/comment/storage"

	"go.uber.org/zap"
)

func Count(ctx context.Context, req *pbcomment.CountReq) (*pbcomment.CountResp, error) {
	// TODO 限制Count

	countMap := make(map[int64]int64, len(req.GetItemIDs()))
	// TODO 这里是不是可以优化SQL 一次给查出来
	for _, v := range req.GetItemIDs() {
		total, err := storage.CountByItemID(ctx, v, req.GetParentIDs()[v])
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		countMap[v] = total
	}
	return &pbcomment.CountResp{
		Total: countMap,
	}, nil
}
