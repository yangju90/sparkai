package vqaservice

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"
	"sparkai/model/qwen"
)

func Service(userId string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = ""
			err = errors.New(fmt.Sprint(r))
		}
	}()

	if v, ok := mem.WSConnContainers[userId]; ok {
		message := v.Messages[len(v.Messages)-1]
		if len(v.ImageData) == 0 {
			return "", errors.New("调用图片问答失败, 没有图片信息")
		}

		data := CreateRequestBody(message, v.ImageData)

		log.Println(v.Messages)
		bytesBody, _ := json.Marshal(data)
		resp, err := http.Post("http://127.0.0.1:8765/chat", "application/json", bytes.NewReader(bytesBody))
		if err != nil {
			log.Println("Error:", err)
			return "", err
		}

		defer resp.Body.Close()

		te := resp.TransferEncoding
		if te != nil && "chunked" == te[0] {
			reader := bufio.NewReader(resp.Body)
			var wsResponse model.WSBodyResponse
			wsResponse.Code = 1
			wsResponse.ContentType = "text"
			wsResponse.Status = 0
			wsResponse.Content = ""

			for {
				chunked, err := reader.ReadString('\n')
				if err != nil {
					return "", err
				}

				var respBody qwen.Qwen2RespBody
				err = json.Unmarshal([]byte(chunked), &respBody)
				if err != nil {
					return "", err
				}

				wsResponse.Content = respBody.Message.Content
				res += wsResponse.Content
				ccc, _ := json.Marshal(wsResponse)
				if e := v.Send(ccc); e != nil {
					return "", e
				}

				if respBody.Done {
					break
				}
			}
		} else {
			return "", errors.New("错误的返回, Not Chunked response!")
		}

		log.Println(res)

		if len(res) != 0 {
			v.AppendMessage(res, constant.ASSISTANT)
		}
	} else {
		return "", errors.New("Id为" + userId + "的用户不在线")
	}

	return res, nil
}

func CreateRequestBody(message model.Message, imageData string) qwen.Qwen2ReqBody {
	request := qwen.Qwen2ReqBody{
		Model: "minicpm",
		Image: imageData,
		Messages: []model.Message{
			message,
		},
	}

	return request
}
