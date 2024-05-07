package io

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sparkai/internal/functionsProcess"
	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"
	"sparkai/model/qwen"
	"strings"

	"github.com/gorilla/websocket"
)

// 0 开始  1 继续文字  2 调用function  9 结束

func Wsservice(userId string) error {
	if v, ok := mem.WSConnContainers[userId]; ok {

		url := "http://192.168.8.232:11434/api/chat"

		body := qwen.NewOllamaReqBody()
		bytesBody, _ := json.Marshal(body)

		resp, err := http.Post(url, "application/json", bytes.NewReader(bytesBody))
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		defer resp.Body.Close()

		te := resp.TransferEncoding
		if te != nil && "chunked" == te[0] {
			fmt.Println("chunked!")
			reader := bufio.NewReader(resp.Body)
			for {
				chunked, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Println("Error:", err)
					break
				}

				var respBody qwen.OllamaRespBody
				err = json.Unmarshal([]byte(chunked), &respBody)
				log.Println(chunked)
			}
		} else {
			return errors.New("错误的返回, Not Chunked response!")
		}
	} else {
		return errors.New("Id为" + userId + "的用户不在线")
	}

	return nil
}

func WaitSparkaiOutput(conn *websocket.Conn, userId string) error {
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

		tmp := ""

		var wsResponse model.WSBodyResponse
		wsResponse.Code = int(code)
		wsResponse.ContentType = "text"
		wsResponse.Status = int(status)
		wsResponse.Content = content

		if wsResponse.Status != 2 {
			answer += content
		} else {
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			fmt.Println("total_tokens:", totalTokens)
			conn.Close()
		}

		if wsResponse.Status == 0 {
			if strings.HasPrefix(wsResponse.Content, "【") {
				tmp += wsResponse.Content
			}
		} else {
			if len(tmp) != 0 {
				tmp += wsResponse.Content
			}
		}

		if wsResponse.Status == 2 || len(tmp) == 0 || len([]rune(tmp)) > 20 {
			funcName := ""
			needCall := false

			if len(tmp) != 0 {
				funcName, needCall = NeedCallFunc(answer)
				wsResponse.Content = tmp
				tmp = ""
			}
			textByteData, err := json.Marshal(wsResponse)
			if err == nil {
				if !needCall {
					if e := v.Send(textByteData); e != nil {
						return e
					}
				}
				wsResponse.Content = ""
				if wsResponse.Status == 2 {
					if needCall {
						funcerr := functionsProcess.ChoiceFuntionCall(funcName, userId)
						if funcerr != nil {
							wsResponse.Content = answer + ",功能调用失败！"
						}
						answer = ""
					}
					wsResponse.Status = 9
					ccc, _ := json.Marshal(wsResponse)
					if e := v.Send(ccc); e != nil {
						return e
					}
				}
			} else {
				return err
			}
		}

		if wsResponse.Status == 9 {
			break
		}
	}

	if len(answer) != 0 {
		v.AppendMessage(answer, constant.ASSISTANT)
	}

	return nil
}

func NeedCallFunc(str string) (string, bool) {
	tmp := []rune(str)
	if string(tmp[0]) != "【" && string(tmp[len(tmp)-1]) != "】" {
		return "", false
	}

	re := regexp.MustCompile(`【(.*?)】`)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 1 {
		log.Println("step 1:" + matches[1])
		funcMsg := strings.Split(matches[1], " ")
		log.Println(len(funcMsg))
		log.Println(funcMsg)
		if len(funcMsg) == 2 {
			if funcMsg[0] == "功能调度" {
				return funcMsg[1], true
			}
		} else {
			return "", false
		}
	}
	return "", false
}
