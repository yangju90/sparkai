package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sparkai/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	interrupt = make(chan os.Signal, 1)
)

var mu sync.Mutex

func main() {
	// 设置信号通道以优雅地关闭程序
	signal.Notify(interrupt, os.Interrupt)

	// WebSocket服务器地址
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8090", Path: "/ws/answer"}
	log.Printf("connecting to %s", u.String())

	// 创建WebSocket客户端连接
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// 启动一个goroutine来接收消息
	go receiveMessages(c)

	// 启动一个goroutine来发送心跳
	go sendHeartbeat(c)

	// 向服务器发送消息
	mu.Lock()

	s := &model.WSBodyRequest{
		Topic: "login",
		ImMessage: model.MessageBody{
			FromId: "1",
		},
	}

	resp, _ := json.Marshal(s)

	err = c.WriteMessage(websocket.TextMessage, resp)
	if err != nil {
		log.Println("write:", err)
		return
	}
	mu.Unlock()

	// 等待用户按下Ctrl+C
	select {
	case <-interrupt:
		// 优雅地关闭WebSocket连接
		c.Close()
	}
}

func receiveMessages(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("收到来自服务器的消息： %s \n", message)
	}
}

func sendHeartbeat(c *websocket.Conn) {
	heartbeatInterval := time.Second * 3 // 设置心跳间隔为30秒
	// pingMsg := []byte("ping")            // 心跳消息内容

	s := &model.WSBodyRequest{
		Topic: "heart_beat",
		ImMessage: model.MessageBody{
			Content: "ping",
		},
	}

	resp, _ := json.Marshal(s)

	for {
		mu.Lock()
		err := c.WriteMessage(websocket.TextMessage, resp)
		mu.Unlock()
		if err != nil {
			log.Println("write:", err)
			return
		}
		time.Sleep(heartbeatInterval) // 等待下一次心跳时间
	}
}
