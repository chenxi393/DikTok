package ws

import (
	"douyin/package/constant"
	"douyin/service"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
)

func HandleWebSocket() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		// c.Locals is added to the *websocket.Conn
		var service service.MessageService
		service.ToUserID, err = strconv.ParseUint(c.Query("to_user_id"), 10, 64)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		userID := c.Locals(constant.UserID).(uint64)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				zap.L().Sugar().Errorf("read:", err)
				break
			}
			zap.L().Sugar().Infof("recv: %s", msg)
			ms := strings.Split(string(msg), "\n")
			if len(ms) == 1 && ms[0] == "get" {
				// 返回聊天记录
				res, err := service.MessageChat(userID)
				if err != nil {
					zap.L().Error(err.Error())
				}
				c.WriteJSON(res)
			} else if len(ms) == 2 && ms[1] == "post" {
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
