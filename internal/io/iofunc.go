package io

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"

	commandIO "sparkai/common/io"
	resources "sparkai/resources"

	"embed"
)

func WaitUserInput(conn *websocket.Conn, appid string) {
	text, err := commandIO.WaitCommandInput()
	if err != nil {
		fmt.Println("Command Input message error:", err)
		text = "你是谁，可以干什么？"
	}
	data := GenParams1(appid, text, true)
	conn.WriteJSON(data)
}

func ReadFunctionsConfig(resources *embed.FS) (configMsg map[string]interface{}, err error) {
	path := "function"
	dir, err := resources.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var regfuncs []map[string]interface{}

	for _, v := range dir {
		vPath := path + "/" + v.Name()
		data, err := resources.ReadFile(vPath)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			break
		}
		var funcMsg map[string]interface{}

		err = json.Unmarshal(data, &funcMsg)
		if err != nil {
			fmt.Println("Error parsing config file:", err)
			break
		}
		regfuncs = append(regfuncs, funcMsg)
	}

	if err == nil {
		configMsg = map[string]interface{}{
			"text": regfuncs,
		}
	}

	return
}

// 生成参数
func GenParams1(appid, question string, first bool) map[string]interface{} {

	messages := []Message{
		{Role: "user", Content: question},
	}

	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": appid,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      "generalv3.5",
				"temperature": float64(0.8),
				"top_k":       int64(6),
				"max_tokens":  int64(2048),
				"auditing":    "default",
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": messages,
			},
		},
	}

	functions, err := ReadFunctionsConfig(&resources.FunctionsConfig)

	if first && err == nil {
		payload := data["payload"].(map[string]interface{})
		payload["functions"] = functions

		fmt.Println(functions)

		fmt.Println("Register function call!")
	}
	return data
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
