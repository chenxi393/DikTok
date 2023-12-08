package handler

import (
	"context"
	"douyin/gateway/auth"
	pbcomment "douyin/grpc/comment"
	"douyin/package/constant"
	"douyin/response"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	CommentClient pbcomment.CommentClient
)

type commentRequest struct {
	ActionType string `query:"action_type"`
	// 要删除的评论id，在action_type=2的时候使用
	CommentID uint64 `query:"comment_id,omitempty"`
	// 用户填写的评论内容，在action_type=1的时候使用
	CommentText *string `query:"comment_text,omitempty"`
	// 视频id
	VideoID uint64 `query:"video_id"`
}

type commentListRequest struct {
	Token string `query:"token"`
	// 视频id
	VideoID uint64 `query:"video_id"`
}

func CommentAction(c *fiber.Ctx) error {
	var req commentRequest
	err := c.QueryParser(&req)
	if err != nil {
		zap.L().Error(err.Error())
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	var resp *pbcomment.CommentResponse
	if req.ActionType == constant.DoAction && req.CommentText != nil {
		resp, err = CommentClient.Add(context.Background(), &pbcomment.AddRequest{
			UserID:  userID,
			VideoID: req.VideoID,
			Content: *req.CommentText,
		})
	} else if req.ActionType == constant.UndoAction && req.CommentID != 0 {
		resp, err = CommentClient.Delete(context.Background(), &pbcomment.DeleteRequest{
			CommentID: req.CommentID,
		})
	} else {
		err = errors.New(response.BadParaRequest)
	}
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
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
			StatusCode:  response.Failed,
			StatusMsg:   response.BadParaRequest,
			CommentList: nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	var userID uint64
	if req.Token == "" {
		userID = 0
	} else {
		claims, err := auth.ParseToken(req.Token)
		if err != nil {
			res := response.CommentListResponse{
				StatusCode:  response.Failed,
				StatusMsg:   response.WrongToken,
				CommentList: nil,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(res)
		}
		userID = claims.UserID
	}
	resp, err := CommentClient.List(context.Background(), &pbcomment.ListRequest{
		UserID:  userID,
		VideoID: req.VideoID,
	})
	if err != nil {
		res := response.CommentActionResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
			Comment:    nil,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
