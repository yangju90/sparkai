package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"

	"github.com/gorilla/websocket"
)

func WaitUserInput(conn *websocket.Conn, appid string, userId string) {
	var messages []model.Message

	if v, ok := mem.WSConnContainers[userId]; ok {
		messages = v.Messages
	} else {
		panic("Id为" + userId + "的用户不在线")
	}

	data := GenParams1(appid, messages, true)

	byteData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("map cast byte error ", err)
		return
	}

	log.Println("发送数据：" + string(byteData))

	conn.WriteMessage(websocket.TextMessage, byteData)
	// conn.WriteJSON(data)
}

func WaitSparkaiOutput(conn *websocket.Conn, userId string) error {
	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			return err
		}

		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			fmt.Println("Error parsing JSON:", err)
			return err1
		}
		fmt.Println(string(msg))
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			fmt.Println(data["payload"])
			return errors.New("sparkai response err")
		}
		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)

		if v, ok := mem.WSConnContainers[userId]; ok {
			textByteData, err := json.Marshal(text)
			if err == nil {
				if e := v.Send(textByteData); e != nil {
					return e
				}
			} else {
				return err
			}
		} else {
			panic("Id为" + userId + "的用户不在线")
		}

		if status != 2 {
			answer += content
		} else {
			fmt.Println("收到最终结果")
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			fmt.Println("total_tokens:", totalTokens)
			conn.Close()
			break
		}

	}
	//输出返回结果
	fmt.Println(answer)

	return nil
}

// 生成参数
func GenParams1(appid string, messages []model.Message, funcCall bool) map[string]interface{} {
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

	if len(constant.FunctionsConfig) != 0 && funcCall {
		payload := data["payload"].(map[string]interface{})
		payload["functions"] = constant.FunctionsConfig
		fmt.Println("Register function call!")
	}
	return data
}
