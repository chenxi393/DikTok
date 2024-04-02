package main

import (
	"context"
	pbcomment "douyin/grpc/comment"
	pbuser "douyin/grpc/user"
	"douyin/model"
	"douyin/package/cache"
	"douyin/package/constant"
	"douyin/package/util"
	"errors"
	"sync"

	"strconv"
	"time"

	"go.uber.org/zap"
)

type CommentService struct {
	pbcomment.UnimplementedCommentServer
}

func (s *CommentService) Add(ctx context.Context, req *pbcomment.AddRequest) (*pbcomment.CommentResponse, error) {
	// TODO 增加敏感词过滤 可以异步实现 comment表多一列屏蔽信息
	id, err := util.GetSonyFlakeID()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 消息队列异步处理评论信息
	msg := &model.Comment{
		ID:          id,
		VideoID:     req.VideoID,
		UserID:      req.UserID,
		Content:     req.Content,
		CreatedTime: time.Now(),
	}
	err = CreateComment(msg)
	if err != nil {
		return nil, err
	}
	// 查找评论的用户信息
	userResponse, err := userClient.Info(ctx, &pbuser.InfoRequest{
		UserID:      req.UserID,
		LoginUserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	commentResponse := &pbcomment.CommentData{
		Id:      msg.ID,
		User:    userResponse.GetUser(),
		Content: msg.Content,
		// 这个评论的时间客户端哈好像可以到毫秒2006-01-02 15:04:05.999
		// 但是感觉每必要 日期就够了
		CreateDate: msg.CreatedTime.Format("2006-01-02 15:04"),
	}
	return &pbcomment.CommentResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.CommentSuccess,
		Comment:    commentResponse,
	}, nil
}

func (s *CommentService) Delete(ctx context.Context, req *pbcomment.DeleteRequest) (*pbcomment.CommentResponse, error) {
	// 我们认为删除评论不是高频动作 故不使用消息队列
	// database里会删缓存 并且校验是不是自己发的 实际上不校验也行
	// 注意还需要在database里减少视频的评论数
	msg, err := DeleteComment(req.CommentID, req.VideoID, req.UserID)
	if err != nil {
		return nil, err
	}
	// 查找评论的用户信息
	userResponse, err := userClient.Info(ctx, &pbuser.InfoRequest{
		UserID:      req.UserID,
		LoginUserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	commentResponse := &pbcomment.CommentData{
		Id:      msg.ID,
		User:    userResponse.GetUser(),
		Content: msg.Content,
		// 这个评论的时间客户端哈好像可以到毫秒2006-01-02 15:04:05.999
		// 但是感觉每必要 日期就够了
		CreateDate: msg.CreatedTime.Format("2006-01-02 15:04"),
	}
	return &pbcomment.CommentResponse{
		StatusCode: constant.Success,
		StatusMsg:  constant.DeleteCommentSuccess,
		Comment:    commentResponse,
	}, nil
}

func (s *CommentService) List(ctx context.Context, req *pbcomment.ListRequest) (*pbcomment.ListResponse, error) {
	// // 使用布隆过滤器判断视频ID是否存在
	// if !VideoIDBloomFilter.TestString(strconv.FormatUint(service.VideoID, 10)) {
	// 	zap.L().Sugar().Error(constant.BloomFilterRejected)
	// 	return nil, fmt.Errorf(constant.BloomFilterRejected)
	// }
	zap.L().Sugar().Infof("[CommentService list] req: %+v", req)
	if req.Count < 0 || req.Count > 50 || req.Offset < 0 {
		return nil, errors.New(constant.BadParaRequest)
	}
	if req.Count == 0 {
		req.Count = 50
	}
	var comments []*model.Comment
	var err error
	if req.GetOffset()+req.Count >= 50 {
		comments, err = GetCommentsByVideoIDFromMaster(req.VideoID, int(req.GetOffset()), int(req.Count))
		if err != nil {
			zap.L().Sugar().Error(err)
			return nil, err
		}
	} else {
		// redis只存50条 评论 多的 去数据库里拿
		comments, err = GetCommentsByVideoID(req.VideoID)
		if err != nil {
			zap.L().Sugar().Warn(constant.CacheMiss)
			// 加分布式锁 这里分布式锁严格测试过了 感觉没什么很大问题
			key := "lock:" + constant.CommentPrefix + strconv.FormatUint(req.VideoID, 10)
			uuidValue, err := util.GetUUid()
			if err != nil {
				return nil, err
			}
			ok := true
			for ok, err = cache.GetLock(key, uuidValue, constant.LockTime, commentRedis); err == nil; {
				if ok {
					defer cache.ReleaseLock(key, uuidValue, videoRedis)
					// redis 只存 50条
					comments, err = GetCommentsByVideoIDFromMaster(req.VideoID, 0, 50)
					if err != nil {
						zap.L().Sugar().Error(err)
						return nil, err
					}
					err := SetComments(req.VideoID, comments)
					if err != nil {
						zap.L().Error(err.Error())
					}
					break
				}
				time.Sleep(constant.RetryTime * time.Millisecond)
				comments, err = GetCommentsByVideoID(req.VideoID)
				if err == nil {
					break
				}
			}
			if err != nil {
				zap.L().Sugar().Error(err)
				return nil, err
			}
		}
		if int(req.GetOffset()+req.Count) <= len(comments) {
			comments = comments[req.GetOffset() : req.GetOffset()+req.Count]
		}
	}

	// 先用map 减少rpc查询次数
	userMap := make(map[uint64]*pbuser.UserInfo)
	for i := range comments {
		userMap[comments[i].UserID] = &pbuser.UserInfo{}
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(userMap))
	for userID := range userMap {
		go func(id uint64) {
			defer wg.Done()
			// TODO 这里是不是也应该 rpc批量拿出来 而不是一个个去拿
			user, err := userClient.Info(ctx, &pbuser.InfoRequest{
				LoginUserID: req.UserID,
				UserID:      id,
			})
			if err != nil {
				zap.L().Error(err.Error())
				// 没有数据 return也是可以接收的 可以重试两次 （封装一下）
				return // 这里 不return 会panic 下面 user
			}
			if user.StatusCode != 0 {
				zap.L().Error("rpc 调用错误")
			}
			// 这里map会不会有并发问题啊
			// TODO 去测试一下
			// 这里如果不用 指针写入的化 会导致下面 videoInfo
			// append 地址被改变 要不就上锁 所有rpc请求之后 再下一个
			// 但是这里之间 直接使用* 似乎也不太好 报了warning
			// 说内部有锁  不能复制 TODO
			*userMap[id] = *user.GetUser()
		}(userID)
	}
	commentResponse := make([]*pbcomment.CommentData, len(comments))
	for i := range comments {
		commentResponse[i] = &pbcomment.CommentData{
			Id:         comments[i].ID,
			User:       userMap[comments[i].UserID],
			Content:    comments[i].Content,
			CreateDate: comments[i].CreatedTime.Format("2006-01-02 15:04"),
		}
	}
	wg.Wait()
	total, err := GetCommentsNumByVideoIDFromMaster(req.VideoID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	hasMore := total > int64(int(req.Offset)+len(commentResponse))
	return &pbcomment.ListResponse{
		StatusCode:  constant.Success,
		StatusMsg:   constant.LoadCommentsSuccess,
		CommentList: commentResponse,
		HasMore:     hasMore,
		Total:       total,
	}, nil
}
