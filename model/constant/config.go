package constant

import (
	"fmt"
	"log"
	"sparkai/model/wssconfig"
	resources "sparkai/resources"

	"gopkg.in/yaml.v2"
)

var WssConfig wssconfig.WssConfig

func init() {
	log.Println("init WssConfig ......")

	var err error
	data, err := resources.ConfigResource.ReadFile("config/application.yaml")

	if err != nil {
		fmt.Println("读取配置文件时发生错误:", err)
		return
	}

	err = yaml.Unmarshal(data, &WssConfig)
	if err != nil {
		fmt.Println("解析配置文件时发生错误:", err)
		return
	}

	fmt.Println(WssConfig.ApiKey)
}
