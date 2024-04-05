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

// 0 开始  1 继续文字  2 调用function  3 结束

func WaitUserInput(conn *websocket.Conn, appid string, userId string) {
	if v, ok := mem.WSConnContainers[userId]; ok {
		data := GenParams1(appid, userId, v.ChatId, v.Messages, v.IsRegistry)

		byteData, err := json.Marshal(data)
		if err != nil {
			log.Println("map cast byte error ", err)
			return
		}

		log.Println("发送数据：" + string(byteData))
		conn.WriteMessage(websocket.TextMessage, byteData)

		// v.IsRegistry = true
	} else {
		panic("Id为" + userId + "的用户不在线")
	}

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
		log.Println(string(msg))
		//解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			log.Println(data["payload"])
			return errors.New("sparkai response err")
		}
		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		functionCall := text[0].(map[string]interface{})["function_call"]

		var wsResponse model.WSBodyResponse
		wsResponse.Code = int(code)

		if functionCall == nil {
			wsResponse.ContentType = "text"
			wsResponse.Status = responseConvert("text", int(status))
			wsResponse.Content = content
		} else {
			functionCallMap := functionCall.(map[string]interface{})
			wsResponse.Status = responseConvert("function", int(status))
			wsResponse.ContentType = "function"
			wsResponse.Content = functionCallMap["name"].(string)
			// Todo function 后续调用
		}

		if v, ok := mem.WSConnContainers[userId]; ok {
			textByteData, err := json.Marshal(wsResponse)
			if err == nil {
				if e := v.Send(textByteData); e != nil {
					return e
				}

				// undo 临时测试
				if wsResponse.Status == 2 {
					wsResponse.Status = 3
					wsResponse.Content += "，调用完成！"
					ccc, _ := json.Marshal(wsResponse)
					if e := v.Send(ccc); e != nil {
						return e
					}
				}

				// undo
			} else {
				return err
			}
		} else {
			panic("Id为" + userId + "的用户不在线")
		}

		if status != 2 {
			answer += content
		} else {
			log.Println("收到最终结果")
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

func responseConvert(contentType string, code int) int {
	var res int = 3
	if contentType == "text" {
		switch code {
		case 0:
			res = 0
		case 1:
			res = 1
		case 2:
			res = 3
		}
	} else if contentType == "function" {
		return code
	}
	return res
}

// 生成参数
func GenParams1(appid string, uid string, chat_id string, messages []model.Message, isRegistry bool) map[string]interface{} {
	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": appid,
			"uid":    uid,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      "generalv3.5",
				"temperature": float64(0.5),
				"top_k":       int64(4),
				"max_tokens":  int64(2048),
				"chat_id":     chat_id,
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": messages,
			},
		},
	}

	if len(constant.FunctionsConfig) != 0 && !isRegistry {
		payload := data["payload"].(map[string]interface{})
		payload["functions"] = constant.FunctionsConfig
		fmt.Println("Register function call!")
	}
	return data
}
