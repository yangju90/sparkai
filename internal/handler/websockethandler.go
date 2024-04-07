package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"
	"strings"
	"sync"

	"github.com/google/uuid"
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
		var requestbody model.WSBodyRequest
		var sessionId string

		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				log.Println("ReadMessage error:", err)
				if v, ok := mem.WSConnContainers[sessionId]; ok {
					v.Close()
				} else {
					conn.Close()
				}
				delete(mem.WSConnContainers, sessionId)
				break
			}

			err = json.Unmarshal(msg, &requestbody)
			var mu sync.Mutex

			if err == nil {
				switch requestbody.Topic {
				case "login":

					messages := []model.Message{
						{
							Role:    constant.SYSTEM,
							Content: constant.FuncPromptConfig,
						},
					}

					mem.WSConnContainers[requestbody.ImMessage.FromId] = &model.WSConnContainer{
						WSConn:     conn,
						MU:         &mu,
						Messages:   messages,
						Status:     "UP",
						ChatId:     strings.ReplaceAll(uuid.New().String(), "-", ""),
						IsRegistry: true,
					}
					sessionId = requestbody.ImMessage.FromId
					if len(sessionId) == 0 {
						conn.Close()
						// break
					} else {
						if v, ok := mem.WSConnContainers[requestbody.ImMessage.FromId]; ok {
							requestbody.ImMessage.Content = "登录成功!"
							responseByte, _ := json.Marshal(requestbody)
							v.Send(responseByte)
						}
					}
				case "logout":
					if v, ok := mem.WSConnContainers[requestbody.ImMessage.FromId]; ok {
						requestbody.ImMessage.Content = "登出成功!"
						responseByte, _ := json.Marshal(requestbody)
						v.Send(responseByte)
						v.Close()
					}
					delete(mem.WSConnContainers, requestbody.ImMessage.FromId)
				case "heart_beat":
					// log.Println(sessionId + "   " + string(msg))

				default:
					mu.Lock()
					err = conn.WriteMessage(websocket.TextMessage, []byte(""))
					mu.Unlock()
					if err != nil {
						log.Println("WriteMessage:", err)
					}
				}
			}
		}
	}()
}
