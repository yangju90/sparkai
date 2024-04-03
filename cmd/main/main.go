package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sparkai/common/gsd"
	"sparkai/internal/sparkaiservice"
	"sparkai/model/wssconfig"

	commandIO "sparkai/common/io"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"
)

// WebSocket 配置文件信息
// Sparkai 每次发送完信息后，服务器会关闭连接
func main() {
	// 解析配置文件
	var wssConfig wssconfig.WssConfig
	readAppConfig(&wssConfig)

	gsd.GracefulShutdown("sparkai")

}

func readAppConfig(wssConfig *wssconfig.WssConfig) {
	data, err := ioutil.ReadFile("application.yaml")
	if err != nil {
		fmt.Println("读取配置文件时发生错误:", err)
		return
	}

	err = yaml.Unmarshal(data, wssConfig)
	if err != nil {
		fmt.Println("解析配置文件时发生错误:", err)
		return
	}
	fmt.Println(wssConfig.ApiKey)
}
