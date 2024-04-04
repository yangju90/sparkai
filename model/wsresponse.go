package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnContainer struct {
	WSConn *websocket.Conn
	MU     *sync.Mutex
	Status string
}

func (container *WSConnContainer) Send(data []byte) {
	container.MU.Lock()
	defer container.MU.Unlock()
	container.WSConn.WriteMessage(websocket.TextMessage, data)
}

func (container *WSConnContainer) Close() error {
	mu := container.MU
	mu.Lock()
	defer mu.Unlock()
	container.MU = nil
	container.Status = "DOWN"
	err := container.WSConn.Close()
	return err
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

func UserIdNotOnline(userId string) *HttpBodyResponse {
	res := HttpBodyResponse{
		Code: 400,
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
