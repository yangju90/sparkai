package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sparkai/common/gsd"
	"sparkai/internal/sparkaiservice"
	"sparkai/model/wssconfig"

	commandIO "sparkai/common/io"

	"gopkg.in/yaml.v2"
)

// WebSocket 配置文件信息
// Sparkai 每次发送完信息后，服务器会关闭连接
func main() {
	// 解析配置文件
	var wssConfig wssconfig.WssConfig
	readAppConfig(&wssConfig)

	sessionId := "1"

	c := make(chan int)
	defer func() {
		c <- 0
	}()

	go func() {
		for {
			text, err := commandIO.WaitCommandInput()
			if err != nil {
				log.Println("Command Input message error:", err)
			}
			select {
			case <-c:
				log.Println("param exit!")
				return
			default:
				sparkaiservice.Wsservice(sessionId, text)
			}

		}
	}()

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
