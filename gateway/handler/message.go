package handler

import (
	"context"
	"douyin/gateway/response"
	pbmessage "douyin/grpc/message"
	"douyin/package/constant"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
	ToUserID int64 `query:"to_user_id"`
}
type messageListRequest struct {
	// 对方用户id
	ToUserID int64 `query:"to_user_id"`
	// //上次最新消息的时间（新增字段-apk更新中）
	Pre_msg_time int64 `query:"pre_msg_time"`
}

func MessageAction(c *fiber.Ctx) error {
	var req sendRequest
	err := c.QueryParser(&req)
	if err != nil {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	if req.ActionType != constant.DoAction {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
			StatusMsg:  constant.BadParaRequest,
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	resp, err := MessageClinet.Send(c.UserContext(), &pbmessage.SendRequest{
		UserID:   userID,
		ToUserID: req.ToUserID,
		Content:  req.Content,
	})
	if err != nil {
		res := response.CommonResponse{
			StatusCode: constant.Failed,
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
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	userID := c.Locals(constant.UserID).(int64)
	resp, err := MessageClinet.List(c.UserContext(), &pbmessage.ListRequest{
		UserID:     userID,
		ToUserID:   req.ToUserID,
		PreMsgTime: req.Pre_msg_time,
	})
	if err != nil {
		res := response.MessageResponse{
			StatusCode: constant.Failed,
			StatusMsg:  err.Error(),
		}
		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
	c.Status(fiber.StatusOK)
	return c.JSON(resp)
}

func MessageWebsocket() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		// c.Locals is added to the *websocket.Conn
		toUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		userID := c.Locals(constant.UserID).(int64)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				zap.L().Sugar().Errorf("read:", err)
				break
			}
			zap.L().Sugar().Infof("recv: %s", msg)
			ms := strings.Split(string(msg), "\n")
			if len(ms) == 1 && ms[0] == "get" {
				// 返回聊天记录
				res, err := MessageClinet.List(context.Background(), &pbmessage.ListRequest{
					UserID:   userID,
					ToUserID: toUserID,
				})
				if err != nil {
					zap.L().Error(err.Error())
				}
				c.WriteJSON(res)
			} else if len(ms) == 2 && ms[1] == "post" {
				// TODO
				// 写入数据库 发送给对应的在线好友
				// 不在线怎么办 ？？
			} else {
				msg = []byte("错误的消息格式 连接关闭")
				err := errors.New(string(msg))
				zap.L().Error(err.Error())
				if err = c.WriteMessage(mt, msg); err != nil {
					zap.L().Sugar().Error("write:", err)
				}
				break
			}
		}
	}
}
