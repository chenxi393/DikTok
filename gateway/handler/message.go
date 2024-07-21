package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"diktok/gateway/response"
	pbmessage "diktok/grpc/message"
	"diktok/package/constant"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
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

type RequestMsg struct {
	Method   string `json:"method"`
	ToUserID int64  `json:"to_user_id"`
	Content  string `query:"content"`
}

func MessageWebsocket() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			msg []byte
			err error
		)
		// c.Locals is added to the *websocket.Conn
		userID := c.Locals(constant.UserID).(int64)
		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				zap.L().Sugar().Errorf("read:", err)
				break
			}
			zap.L().Sugar().Infof("recv: %s", msg)
			var msgJson RequestMsg
			err := json.Unmarshal(msg, &msgJson)
			if err != nil {
				err = c.WriteJSON(&response.CommonResponse{
					StatusCode: -1,
					StatusMsg:  constant.BadParaRequest,
				})
				if err != nil {
					zap.L().Sugar().Error("write:", err)
					break
				}
				continue
			}
			zap.L().Sugar().Infof("recvJSON: %#v", msgJson)
			if msgJson.Method == "get" {
				// 返回聊天记录
				res, err := MessageClinet.List(context.Background(), &pbmessage.ListRequest{
					UserID:   userID,
					ToUserID: msgJson.ToUserID,
				})
				if err != nil {
					zap.L().Error(err.Error())
				}
				c.WriteJSON(res)
			} else if msgJson.Method == "post" {
				// TODO
				// 写入数据库 发送给对应的在线好友
				// 不在线怎么办 ？？
			} else {
				err = c.WriteJSON(&response.CommonResponse{
					StatusCode: -1,
					StatusMsg:  constant.BadParaRequest,
				})
				if err != nil {
					zap.L().Sugar().Error("write:", err)
					break
				}
			}
			// TODO 存储 在线状态
			// 当用户收到消息的时候 主动推送
		}
	}
}

func SSEHandle(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		var msgJson RequestMsg
		err := c.BodyParser(&msgJson)
		if err != nil {
			err = c.JSON(&response.CommonResponse{
				StatusCode: -1,
				StatusMsg:  constant.BadParaRequest,
			})
			if err != nil {
				zap.L().Sugar().Error(err.Error())
			}
			return
		}
		userID := c.Locals(constant.UserID).(int64)
		switch msgJson.Method {
		case "chat":
			{
				stream, err := MessageClinet.RequestToLLM(c.Context(), &pbmessage.RequestToLLMRequest{
					UserID:  userID,
					Content: msgJson.Content,
				})
				if err != nil {
					zap.L().Sugar().Error(err.Error())
				}

				// 3. for循环获取服务端推送的消息
				for {
					// 通过 Recv() 不断获取服务端send()推送的消息
					resp, err := stream.Recv()
					// 4. err==io.EOF则表示服务端关闭stream了 退出
					if err == io.EOF {
						zap.L().Sugar().Infof("server closed")
						break
					}
					if err != nil {
						zap.L().Sugar().Errorf("Recv error:%v", err)
						break
					}
					zap.L().Sugar().Infof("Recv data:%v", resp.String())
					fmt.Fprintf(w, "data: %s\n\n", resp.String())
					if err := w.Flush(); err != nil {
						break
					}
				}
			}
		}
	}))
	return nil
}
