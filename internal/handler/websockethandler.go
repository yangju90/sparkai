package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域请求
	},
}

func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// 处理WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	// defer conn.Close()

	conn.SetPingHandler(func(message string) error {
		fmt.Println(message)
		err := conn.WriteMessage(websocket.PongMessage, []byte("pong"))
		return err
	})

	go func() {
		for {
			// 读取消息

			mt, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage error:", err)
				break
			}

			fmt.Printf("收到消息：%s\n", string(msg))

			fmt.Println(mt)

			// 检查消息类型
			// switch msg := msg.(type) {
			// case *websocket.PingMessage:
			// 	// 发送 PongMessage 作为响应
			// 	err := conn.WriteMessage(websocket.PongMessage, msg.Data)
			// 	if err != nil {
			// 		log.Println("WriteMessage error:", err)
			// 		break
			// 	}
			// }

			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("WriteMessage:", err)
				break
			}
		}
	}()
}
