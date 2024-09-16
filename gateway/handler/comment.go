package handler

import (
	"errors"

	"diktok/gateway/response"
	pbcomment "diktok/grpc/comment"
	"diktok/package/constant"
	"diktok/package/rpc"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type commentRequest struct {
	ActionType  string `query:"action_type"` // 要删除的评论id，在action_type=2的时候使用
	CommentID   int64  `query:"comment_id,omitempty"`
	CommentText string `query:"comment_text,omitempty"` // 用户填写的评论内容，在action_type=1的时候使用
	VideoID     int64  `query:"video_id"`               // 视频id
	ParentID    int64  `query:"parent_id"`
}

type commentListRequest struct {
	VideoID       int64 `query:"video_id"` // 视频id
	Count         int32 `query:"count"`
	LastCommentId int64 `query:"last_comment_id"`
}

func CommentAction(c *fiber.Ctx) error {
	var req commentRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommentActionResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	var resp *pbcomment.CommentResponse
	if req.ActionType == constant.DoAction && req.CommentText != "" {
		resp, err = rpc.CommentClient.Add(c.UserContext(), &pbcomment.AddRequest{
			UserID:   userID,
			VideoID:  req.VideoID,
			Content:  req.CommentText,
			ParentID: req.ParentID,
		})
	} else if req.ActionType == constant.UndoAction && req.CommentID != 0 {
		resp, err = rpc.CommentClient.Delete(c.UserContext(), &pbcomment.DeleteRequest{
			VideoID:   req.VideoID,
			CommentID: req.CommentID,
			UserID:    userID,
		})
	} else {
		err = errors.New(constant.BadParaRequest)
	}
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func CommentList(c *fiber.Ctx) error {
	var req commentListRequest
	err := c.QueryParser(&req)
	if err != nil {
		res := response.CommentListResponse{
			StatusCode:  constant.Failed,
			StatusMsg:   constant.BadParaRequest,
			CommentList: nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	resp, err := rpc.CommentClient.List(c.UserContext(), &pbcomment.ListRequest{
		UserID:        userID,
		VideoID:       req.VideoID,
		LastCommentId: req.LastCommentId,
		Count:         req.Count,
	})
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
