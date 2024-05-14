package io

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sparkai/internal/functionsProcess"
	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"
	"sparkai/model/qwen"
	"strings"
)

// 0 开始  1 继续文字  2 调用function  9 结束

func Wsservice(userId string) error {
	if v, ok := mem.WSConnContainers[userId]; ok {

		url := "http://192.168.8.232:11434/api/chat"

		var answer = ""

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
			var wsResponse model.WSBodyResponse
			wsResponse.Code = 0
			wsResponse.ContentType = "text"
			wsResponse.Status = 0
			wsResponse.Content = ""

			tmp := ""
			index := 0
			for {
				chunked, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				var respBody qwen.OllamaRespBody
				err = json.Unmarshal([]byte(chunked), &respBody)
				if err != nil {
					return err
				}

				if index == 0 {
					wsResponse.Status = 0
				} else {
					if respBody.Done {
						wsResponse.Status = 2
					} else {
						wsResponse.Status = 1
					}
				}
				wsResponse.Content = respBody.Message.Content

				// answer 添加， content改写
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

				index++

				if wsResponse.Status == 9 {
					break
				}
			}
		} else {
			return errors.New("错误的返回, Not Chunked response!")
		}

		if len(answer) != 0 {
			v.AppendMessage(answer, constant.ASSISTANT)
		}
	} else {
		return errors.New("Id为" + userId + "的用户不在线")
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
