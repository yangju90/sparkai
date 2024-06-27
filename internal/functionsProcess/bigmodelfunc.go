package functionsProcess

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"
	"sparkai/model/qwen"
)

func BigModelFunc(userId string) (string, error) {
	if v, ok := mem.WSConnContainers[userId]; ok {
		body := qwen.CreateNewQwen2ReqBodyWithGeneralPrompt(v.Messages)
		bytesBody, _ := json.Marshal(body)

		resp, err := http.Post(constant.WssConfig.HostUrl, "application/json", bytes.NewReader(bytesBody))
		if err != nil {
			log.Println("Error:", err)
			return "", err
		}
		defer resp.Body.Close()

		te := resp.TransferEncoding
		if te != nil && "chunked" == te[0] {
			reader := bufio.NewReader(resp.Body)
			var wsResponse model.WSBodyResponse
			wsResponse.Code = 0
			wsResponse.ContentType = "text"
			wsResponse.Status = 1
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

				if respBody.Done {
					wsResponse.Status = 2
				} else {
					wsResponse.Status = 1
				}

				wsResponse.Content = respBody.Message.Content
				textByteData, err := json.Marshal(wsResponse)
				if err == nil {
					if e := v.Send(textByteData); e != nil {
						return "", e
					}
				} else {
					return "", err
				}

				if wsResponse.Status == 2 {
					break
				}
			}
		} else {
			return "", errors.New("错误的返回, Not Chunked response!")
		}
	} else {
		return "", errors.New("Id为" + userId + "的用户不在线")
	}

	return "", nil
}
