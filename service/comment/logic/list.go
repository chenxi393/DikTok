package logic

import (
	"context"
	"errors"

	pbcomment "diktok/grpc/comment"
	"diktok/package/constant"
	"diktok/service/comment/storage"
	"diktok/storage/database"
	"diktok/storage/database/model"
	"diktok/storage/database/query"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// 条件处理 -> 数据加载 -> 处理打包
func List(ctx context.Context, req *pbcomment.ListRequest) (*pbcomment.ListResponse, error) {
	/*-------------条件处理----------------*/
	if req.Limit < 0 || req.Limit > 50 {
		return nil, errors.New(constant.BadParaRequest)
	}

	var orderBy []field.Expr
	var conds []gen.Condition
	var resp = &pbcomment.ListResponse{}

	so := query.Use(database.DB).CommentMetum
	conds = append(conds, so.ItemID.Eq(req.GetItemID()))
	conds = append(conds, so.ParentID.Eq(req.GetParentID()))
	if len(req.Status) != 0 {
		conds = append(conds, so.Status.In(req.GetStatus()...))
	}
	if req.MaxCommentId != 0 {
		conds = append(conds, so.CommentID.Lt(req.GetMaxCommentId()))
	}
	if req.SortType == 1 {
		orderBy = append(orderBy, so.CreatedAt.Desc())
	} else if req.SortType == 2 {
		orderBy = append(orderBy, so.CreatedAt.Asc())
	}

	// 特殊逻辑 如果有评论id 前置条件作废
	// 是不是 再抽出一个接口比较好
	if req.GetCommentID() != 0 {
		conds = []gen.Condition{so.CommentID.Eq(req.GetCommentID())}
	}

	/*-------------数据加载----------------*/
	// 获取评论元信息
	// 注意这里 limit+1 判断是否has_more  count如果不需要就不返回
	// TODO 这一块 缓存看怎么搞
	comments, err := storage.MGetCommentsByCond(ctx, int(req.Offset), int(req.Limit)+1, conds, orderBy...)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}

	resp.HasMore = len(comments) > int(req.GetLimit())
	if resp.HasMore {
		comments = comments[:len(comments)-1]
	}
	if req.NeedTotal {
		resp.Total, err = storage.CountByCond(ctx, conds)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
	}
	// 空则 huireturn
	if len(comments) == 0 {
		resp.CommentList = nil
		return resp, nil
	}

	commentIDs := make([]int64, 0, len(comments))
	for _, v := range comments {
		commentIDs = append(commentIDs, v.CommentID)
	}
	commentsContent, err := storage.GetContentByCache(ctx, commentIDs)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}
	commentContentMp := make(map[int64]*model.CommentContent, len(commentsContent))
	for _, v := range commentsContent {
		commentContentMp[v.ID] = v
	}

	/*-------------处理打包----------------*/
	var extraJson CommentExtra
	commentResponse := make([]*pbcomment.CommentData, 0, len(comments))
	for _, v := range comments {
		cont := commentContentMp[v.CommentID]
		sonic.Unmarshal([]byte(cont.Extra), &extraJson)

		commentResponse = append(commentResponse, &pbcomment.CommentData{
			CommentID: v.CommentID,
			ItemID:    v.ItemID,
			UserID:    v.UserID,
			ParentID:  v.ParentID,
			CreateAt:  v.CreatedAt.Unix(),
			Status:    v.Status,

			Content:     cont.Content,
			ToCommentID: extraJson.ToCommentID,
			ImageURI:    extraJson.ImageURI,
		})
	}

	resp.CommentList = commentResponse
	return resp, nil

	// var comments []*model.Comment
	// var err error
	// if req.GetOffset()+req.Count >= 50 {
	// 	comments, err = GetCommentsByVideoIDFromMaster(req.VideoID, int(req.GetOffset()), int(req.Count))
	// 	if err != nil {
	// 		zap.L().Sugar().Error(err)
	// 		return nil, err
	// 	}
	// } else {
	// 	// redis只存50条 评论 多的 去数据库里拿
	// 	comments, err = GetCommentsByVideoID(req.VideoID)
	// 	if err != nil {
	// 		zap.L().Sugar().Warn(constant.CacheMiss)
	// 		// 加分布式锁 这里分布式锁严格测试过了 感觉没什么很大问题
	// 		key := "lock:" + constant.CommentPrefix + strconv.FormatUint(req.VideoID, 10)
	// 		uuidValue, err := util.GetUUid()
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		ok := true
	// 		for ok, err = cache.GetLock(key, uuidValue, constant.LockTime, commentRedis); err == nil; {
	// 			if ok {
	// 				defer cache.ReleaseLock(key, uuidValue, videoRedis)
	// 				// redis 只存 50条
	// 				comments, err = GetCommentsByVideoIDFromMaster(req.VideoID, 0, 50)
	// 				if err != nil {
	// 					zap.L().Sugar().Error(err)
	// 					return nil, err
	// 				}
	// 				err := SetComments(req.VideoID, comments)
	// 				if err != nil {
	// 					zap.L().Error(err.Error())
	// 				}
	// 				break
	// 			}
	// 			time.Sleep(constant.RetryTime * time.Millisecond)
	// 			comments, err = GetCommentsByVideoID(req.VideoID)
	// 			if err == nil {
	// 				break
	// 			}
	// 		}
	// 		if err != nil {
	// 			zap.L().Sugar().Error(err)
	// 			return nil, err
	// 		}
	// 	}
	// 	if int(req.GetOffset()+req.Count) <= len(comments) {
	// 		comments = comments[req.GetOffset() : req.GetOffset()+req.Count]
	// 	}
	// }

}
