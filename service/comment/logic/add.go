package logic

import (
	"context"

	pbbase "diktok/grpc/base"
	pbcomment "diktok/grpc/comment"
	"diktok/package/constant"
	"diktok/package/util"
	"diktok/service/comment/storage"
	"diktok/storage/database/model"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
)

func Add(ctx context.Context, req *pbcomment.AddRequest) (*pbbase.BaseResp, error) {
	// TODO 增加敏感词过滤 可以异步实现 comment表多一列屏蔽信息
	id, err := util.GetSonyFlakeID()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	var extra *CommentExtra
	var extraString string
	if req.ToCommentID != 0 || req.ImageURI != "" {
		extra = &CommentExtra{
			ToCommentID: req.ToCommentID,
			ImageURI:    req.ImageURI,
		}
		extraString, err = sonic.MarshalString(extra)
		if err != nil {
			return nil, err
		}
	}
	// 这块写入逻辑可以走MQ 甚至1s 批量写入（flink 大数据批量写）
	// 评论元信息表
	meta := &model.CommentMetum{
		CommentID: int64(id),
		ItemID:    req.ItemID,
		UserID:    req.UserID,
		ParentID:  req.ParentID,
	}
	// 评论内容表
	content := &model.CommentContent{
		ID:      int64(id),
		Content: req.Content,
		Extra:   extraString,
	}
	err = storage.CreateCommentContent(ctx, meta, content)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &pbbase.BaseResp{
		StatusCode: constant.Success,
		StatusMsg:  constant.CommentSuccess,
	}, nil
}
