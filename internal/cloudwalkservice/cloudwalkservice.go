package cloudwalkservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sparkai/model/mem"
)

func Service(userId string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = ""
			err = errors.New(fmt.Sprint(r))
		}
	}()
	// 准备请求体数据
	var data map[string]interface{}

	if v, ok := mem.WSConnContainers[userId]; ok {
		message := v.Messages[len(v.Messages)-1]
		if len(v.ImageData) == 0 {
			return "", errors.New("调用图片问答失败, 没有图片信息")
		}
		data = CreateRequestBody(message.Content, v.ImageData)
	} else {
		return "", errors.New("Id为" + userId + "的用户不在线")
	}

	requestData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "https://maastest.cloudwalk.com:24430/api/gpt/large_model", bytes.NewBuffer(requestData))
	if err != nil {
		fmt.Println("创建云从请求失败:", err)
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("maasKey", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJlMTg2ZDZjMzRmMzM0NDQ4OGNjZDc3MDA0MzU2NTgxNCIsInBhcmVudElkIjoiIiwib3JnIjoiIiwianRpIjoiZjIzODkyZTY2MDk5NGM5Yjg1NzQxZTQ0YzRjNGQxNDYifQ.dPOO2oFzABKt_MebWKFeLDoIza0p9jyZaDN8Y5fCkxo") // 添加自定义头部

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送云从请求失败:", err)
		return "", err
	}

	defer resp.Body.Close()

	// 读取响应
	var responseBuf bytes.Buffer
	_, err = responseBuf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return "", err
	}

	var response map[string]interface{}
	err = json.Unmarshal(responseBuf.Bytes(), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", err
	}

	log.Println(response)
	//解析数据
	attrs := response["attrs"].(map[string]interface{})
	choices := attrs["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	return content, nil
}
