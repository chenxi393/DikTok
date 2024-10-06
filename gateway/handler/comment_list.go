package handler

import (
	"diktok/gateway/response"
	pbcomment "diktok/grpc/comment"
	pbuser "diktok/grpc/user"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
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
	for _, v := range resp.GetCommentList() {
		userIDs = append(userIDs, v.UserID)
	}
	userResp, err := rpc.UserClient.List(c.Context(), &pbuser.ListReq{
		UserID:      userIDs,
		LoginUserID: userID,
	})
	if err != nil {
		return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
	}
	return c.JSON(PackCommentList(resp, userResp))
}

func PackCommentList(comList *pbcomment.ListResponse, userList *pbuser.ListResp) *response.CommentListResponse {
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
			res.CommentList = append(res.CommentList, response.BuildComment(v, UserMap))
		}
	}
	return res
}
