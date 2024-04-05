package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnContainer struct {
	WSConn     *websocket.Conn
	MU         *sync.Mutex
	Status     string
	Messages   []Message
	ChatId     string
	IsRegistry bool
}

func (container *WSConnContainer) Send(data []byte) error {
	container.MU.Lock()
	defer container.MU.Unlock()
	return container.WSConn.WriteMessage(websocket.TextMessage, data)
}

func (container *WSConnContainer) Close() error {
	mu := container.MU
	mu.Lock()
	defer mu.Unlock()
	container.MU = nil
	container.Status = "DOWN"
	container.Messages = nil
	err := container.WSConn.Close()
	return err
}

func (container *WSConnContainer) AppendMessage(text string, messageType string) {
	current := Message{
		Role:    messageType,
		Content: text,
	}
	container.Messages = append(container.Messages, current)
}

func (container *WSConnContainer) NewMessages(text string, messageType string) {
	current := []Message{{
		Role:    messageType,
		Content: text,
	},
	}
	container.Messages = current
}

type HttpBodyRequest struct {
	Text      string `json:"text"`
	ImageData string `json:"imageData"`
}

type HttpBodyResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Success() *HttpBodyResponse {
	res := HttpBodyResponse{
		Code: 200,
		Msg:  "成功",
	}
	return &res
}

func Faild(err error) *HttpBodyResponse {
	res := HttpBodyResponse{
		Code: 40001,
		Msg:  "调用失败！" + err.Error(),
	}
	return &res
}

func UserIdNotOnline(userId string) *HttpBodyResponse {
	res := HttpBodyResponse{
		Code: 20001,
		Msg:  "Id为" + userId + "的用户不在线",
	}
	return &res
}

type WSBodyRequest struct {
	Topic     string      `json:"topic"`
	Device    string      `json:"device"`
	ImMessage MessageBody `json:"imMessage"`
}

type MessageBody struct {
	FromId  string `json:"fromId"`
	Content string `json:"content"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type WSBodyResponse struct {
	Code        int    `json:"code"`
	Status      int    `json:"status"`
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}
