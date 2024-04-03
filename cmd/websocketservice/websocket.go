package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域请求
	},
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	// 处理POST请求
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		// 解析JSON数据并进行处理
		// ...
		fmt.Fprintf(w, "收到POST请求")
	} else {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	}
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user/question", handlePostRequest).Methods("POST")
	router.HandleFunc("/ws/answer", handleWebSocketConnection)

	log.Fatal(http.ListenAndServe(":8080", router))
}
