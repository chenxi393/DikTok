package handler

import (
	"diktok/gateway/response"
	pbcomment "diktok/grpc/comment"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"
	"diktok/package/util"
	"sync"
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type commentListRequest struct {
	VideoID  int64 `query:"video_id"` // 视频id
	ParentID int64 `query:"parent_id"`
	SortType int32 `query:"sort_type"`
	Limit    int32 `query:"limit"`
	Offset   int32 `query:"offset"`
}

func CommentList(c *fiber.Ctx) error {
	var req commentListRequest
	err := c.QueryParser(&req)
	if err != nil {
		return c.JSON(constant.InvalidParams)
	}
	userID := c.Locals(constant.UserID).(int64)
	if req.ParentID == 0 {
		req.ParentID = req.VideoID
	}
	resp, err := rpc.CommentClient.List(c.UserContext(), &pbcomment.ListRequest{
		ItemID:    req.VideoID,
		ParentID:  req.ParentID,
		SortType:  req.SortType,
		NeedTotal: true,
		Offset:    req.Offset,
		Limit:     req.Limit,
	})
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}

	// 查找评论的用户信息
	userIDs := make([]int64, 0, len(resp.CommentList))
	parentsIDs := make([]int64, 0, len(resp.CommentList)) // 查找子评论
	for _, v := range resp.GetCommentList() {
		userIDs = append(userIDs, v.UserID)
		// 父评论才去拉取
		if v.ItemID == v.ParentID {
			parentsIDs = append(parentsIDs, v.CommentID)
		}
	}
	zap.L().Sugar().Debugf("[CommentList] userIDs = %s, parentsIDs= %s", util.GetLogStr(userIDs), util.GetLogStr(parentsIDs))
	wg := sync.WaitGroup{}
	wg.Add(2)
	var userResp *pbuser.ListResp
	var countResp *pbcomment.CountResp
	var errCount int32
	go func() {
		defer wg.Done()
		var err error
		userResp, err = rpc.UserClient.List(c.Context(), &pbuser.ListReq{
			UserID:      userIDs,
			LoginUserID: userID,
		})
		if err != nil {
			atomic.AddInt32(&errCount, 1)
		}
	}()
	go func() {
		defer wg.Done()
		if len(parentsIDs) == 0 {
			return
		}
		var err error
		countResp, err = rpc.CommentClient.Count(c.Context(), &pbcomment.CountReq{
			ParentIDs:   parentsIDs,
			ItemIdIndex: req.VideoID, // 这里是子评论
		})
		if err != nil {
			atomic.AddInt32(&errCount, 1)
		}
	}()
	wg.Wait()
	if errCount > 0 {
		return c.JSON(constant.ServerInternal.WithDetails("errCount > 0"))
	}
	return c.JSON(PackCommentList(resp, userResp, countResp))
}

func PackCommentList(comList *pbcomment.ListResponse, userList *pbuser.ListResp, countResp *pbcomment.CountResp) *response.CommentListResponse {
	res := &response.CommentListResponse{
		CommentList: make([]*response.Comment, 0, len(comList.GetCommentList())),
		HasMore:     comList.GetHasMore(),
		Total:       comList.GetTotal(),
		StatusCode:  constant.Success,
		StatusMsg:   constant.LoadCommentsSuccess,
	}

	UserMap := response.BuildUserMap(userList)
	for _, v := range comList.GetCommentList() {
		if v != nil {
			res.CommentList = append(res.CommentList, response.BuildComment(v, UserMap, countResp.GetCountMap()))
		}
	}
	return res
}
