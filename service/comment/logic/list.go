package logic

import (
	"context"
	"errors"

	pbcomment "diktok/grpc/comment"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"
	"diktok/service/comment/storage"
	"diktok/storage/database/model"

	"go.uber.org/zap"
)

func List(ctx context.Context, req *pbcomment.ListRequest) (*pbcomment.ListResponse, error) {
	// // 使用布隆过滤器判断视频ID是否存在
	// if !VideoIDBloomFilter.TestString(strconv.FormatUint(service.VideoID, 10)) {
	// 	zap.L().Sugar().Error(constant.BloomFilterRejected)
	// 	return nil, fmt.Errorf(constant.BloomFilterRejected)
	// }
	zap.L().Sugar().Infof("[CommentService list] req: %+v", req)
	if req.Count < 0 || req.Count > 50 {
		return nil, errors.New(constant.BadParaRequest)
	}
	if req.Count == 0 {
		req.Count = 50
	}
	if req.LastCommentId == 0 {
		id, err := util.GetSonyFlakeID()
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		req.LastCommentId = int64(id)
	}
	var comments []*model.Comment
	var err error
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

	comments, err = storage.GetCommentsByVideoIDFromMaster(req.VideoID, req.GetLastCommentId(), req.GetCount()+1)
	if err != nil {
		zap.L().Sugar().Error(err)
		return nil, err
	}
	hasMore := len(comments) > int(req.GetCount())
	if hasMore {
		comments = comments[:len(comments)-1]
	}
	// 先用map 减少rpc查询次数
	userIDs := make([]int64, 0)
	for _, c := range comments {
		userIDs = append(userIDs, c.UserID)
	}

	userResp, err := rpc.UserClient.List(ctx, &pbuser.ListReq{
		UserID:      userIDs,
		LoginUserID: req.UserID,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	commentResponse := make([]*pbcomment.CommentData, len(comments))
	for i := range comments {
		commentResponse[i] = &pbcomment.CommentData{
			Id:         comments[i].ID,
			User:       userResp.User[comments[i].UserID],
			Content:    comments[i].Content,
			CreateDate: comments[i].CreatedTime.Format("2006-01-02 15:04"),
		}
	}

	total, err := storage.GetCommentsNumByVideoIDFromMaster(req.VideoID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &pbcomment.ListResponse{
		StatusCode:  constant.Success,
		StatusMsg:   constant.LoadCommentsSuccess,
		CommentList: commentResponse,
		HasMore:     hasMore,
		Total:       total,
	}, nil
}
