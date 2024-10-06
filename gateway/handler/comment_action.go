package handler

import (
	pbcomment "diktok/grpc/comment"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type commentRequest struct {
	ActionType  string `query:"action_type"`
	CommentID   int64  `query:"comment_id"`   // 要删除的评论id，在action_type=2的时候使用
	CommentText string `query:"comment_text"` // 用户填写的评论内容，在action_type=1的时候使用
	VideoID     int64  `query:"video_id"`     // 视频id
	ParentID    int64  `query:"parent_id"`
	ImageUrI    string `query:"image_uri"`
	ToCommentID int64  `query:"to_comment_id"`
}

func CommentAction(c *fiber.Ctx) error {
	var req commentRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		return c.JSON(constant.InvalidParams)
	}
	userID := c.Locals(constant.UserID).(int64)
	if req.ActionType == constant.DoAction && req.CommentText != "" {
		// TODO 需要先查询视频id 是否存在
		resp, err := rpc.CommentClient.Add(c.UserContext(), &pbcomment.AddRequest{
			ItemID:      req.VideoID,
			UserID:      userID,
			Content:     req.CommentText,
			ParentID:    req.ParentID,
			ImageURI:    req.ImageUrI,
			ToCommentID: req.ToCommentID,
		})
		if err != nil {
			return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
		}
		if resp.StatusCode != 0 {
			return c.JSON(constant.ServerLogic.WithDetails(resp.StatusMsg))
		}
		return c.JSON(constant.ServerSuccess)
	} else if req.ActionType == constant.UndoAction && req.CommentID != 0 {
		ListResp, err := rpc.CommentClient.List(c.UserContext(), &pbcomment.ListRequest{
			CommentID: req.CommentID, // 要删除的ID
			Limit:     1,
			Offset:    0,
		})
		if err != nil {
			return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
		}
		if len(ListResp.CommentList) <= 0 {
			return c.JSON(constant.ItemNotFound.WithDetails("评论不存在"))
		}
		commentInfo := ListResp.CommentList[0]
		if commentInfo.GetUserID() != userID {
			return c.JSON(constant.IllegalOperation.WithDetails("不能删除别人的评论"))
		}
		resp, err := rpc.CommentClient.Delete(c.UserContext(), &pbcomment.DeleteRequest{
			CommentID: req.CommentID,
		})
		if err != nil {
			return c.JSON(constant.ServerInternal.WithDetails(err.Error()))
		}
		if resp.StatusCode != 0 {
			return c.JSON(constant.ServerLogic.WithDetails(resp.StatusMsg))
		}
		return c.JSON(constant.ServerSuccess)
	}
	// 最后返回错误
	return c.JSON(constant.InvalidParams)
}
