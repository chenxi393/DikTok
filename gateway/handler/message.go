package handler

import (
	"context"
	pbmessage "douyin/grpc/message"
	"douyin/package/constant"
	"douyin/response"

	"github.com/gofiber/fiber/v2"
)

var (
	MessageClinet pbmessage.MessageClient
)

type sendRequest struct {
	// 1-发送消息
	ActionType string `query:"action_type"`
	// 消息内容
	Content string `query:"content"`
	// 对方用户id
	ToUserID uint64 `query:"to_user_id"`
}
type messageListRequest struct {
	// 对方用户id
	ToUserID uint64 `query:"to_user_id"`
	// //上次最新消息的时间（新增字段-apk更新中）
	Pre_msg_time int64 `query:"pre_msg_time"`
}

func MessageAction(c *fiber.Ctx) error {
	var req sendRequest
	err := c.QueryParser(&req)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	if req.ActionType != constant.DoAction {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  response.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	resp, err := MessageClinet.Send(context.Background(), &pbmessage.SendRequest{
		UserID:  userID,
		ToUerID: req.ToUserID,
		Content: req.Content,
	})
	if err != nil {
		res := response.CommonResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

// 一旦进入消息界面 客户端每秒会调用一次
func MessageChat(c *fiber.Ctx) error {
	var req messageListRequest
	err := c.QueryParser(&req)
	if err != nil {
		res := response.MessageResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(uint64)
	resp, err := MessageClinet.List(context.Background(), &pbmessage.ListRequest{
		UserID:     userID,
		ToUerID:    req.ToUserID,
		PreMsgTime: req.Pre_msg_time,
	})
	if err != nil {
		res := response.MessageResponse{
			StatusCode: response.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}
