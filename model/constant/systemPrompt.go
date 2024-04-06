package constant

import (
	"log"
	"os"
	resources "sparkai/resources"
)

var SystemPromptConfig string

func init() {

	path := "E:/goconfig/system/systemPromptConfig.json"

	_, err := os.Stat(path)
	if err != nil {
		log.Println("local system prompt not exists: " + path)
		data, err := resources.SystemConfig.ReadFile("system/systemPromptConfig.json")
		if err != nil {
			log.Println("read system prompt info error:", err)
			return
		}
		SystemPromptConfig = string(data)
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Println("read file error:", err)
			return
		}
		SystemPromptConfig = string(data)
	}
}
