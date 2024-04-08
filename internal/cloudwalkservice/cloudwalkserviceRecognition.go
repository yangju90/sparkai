package cloudwalkservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sparkai/model/mem"
)

func ServiceRecognition(userId string) (res []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = errors.New(fmt.Sprint(r))
		}
	}()
	// 准备请求体数据
	var data map[string]interface{}

	if v, ok := mem.WSConnContainers[userId]; ok {
		if len(v.ImageData) == 0 {
			return nil, errors.New("调用图片识别失败, 没有图片信息")
		}
		data = CreateRequestBody("人,车,装备", v.ImageData, "ldm")
	} else {
		return nil, errors.New("Id为" + userId + "的用户不在线")
	}

	requestData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "https://maastest.cloudwalk.com:24430/api/gpt/large_model", bytes.NewBuffer(requestData))
	if err != nil {
		fmt.Println("创建云从请求失败:", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("maasKey", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJlMTg2ZDZjMzRmMzM0NDQ4OGNjZDc3MDA0MzU2NTgxNCIsInBhcmVudElkIjoiIiwib3JnIjoiIiwianRpIjoiZjIzODkyZTY2MDk5NGM5Yjg1NzQxZTQ0YzRjNGQxNDYifQ.dPOO2oFzABKt_MebWKFeLDoIza0p9jyZaDN8Y5fCkxo") // 添加自定义头部

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送云从请求失败:", err)
		return nil, err
	}

	defer resp.Body.Close()

	// 读取响应
	var responseBuf bytes.Buffer
	_, err = responseBuf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(responseBuf.Bytes(), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}
	attrs := response["attrs"].(map[string]interface{})
	result := attrs["result"].([]interface{})

	return result, nil
}
