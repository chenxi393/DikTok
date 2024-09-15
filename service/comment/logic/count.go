package logic

import (
	"context"

	pbcomment "diktok/grpc/comment"
	"diktok/service/comment/storage"

	"go.uber.org/zap"
)

func Count(ctx context.Context, req *pbcomment.CountReq) (*pbcomment.CountResp, error) {
	countMap := make(map[int64]int64, len(req.GetVideoID()))
	for _, v := range req.GetVideoID() {
		total, err := storage.GetCommentsNumByVideoIDFromMaster(v)
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
