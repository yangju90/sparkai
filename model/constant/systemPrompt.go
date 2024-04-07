package constant

import (
	"log"
	"os"
	resources "sparkai/resources"
)

var FuncPromptConfig string
var GeneralPromptConfig string

func init() {
	readSystemPrompt("funcPromptConfig.prompt", &FuncPromptConfig)
	// log.Println(FuncPromptConfig)
	readSystemPrompt("generalPromptConfig.prompt", &GeneralPromptConfig)
	// log.Println(GeneralPromptConfig)
}

func readSystemPrompt(fileName string, content *string) {
	path := "E:/goconfig/system/" + fileName

	_, err := os.Stat(path)
	if err != nil {
		log.Println("local system prompt not exists: " + path)
		data, err := resources.SystemConfig.ReadFile("system/" + fileName)
		if err != nil {
			log.Println("read system prompt info error:", err)
			return
		}
		*content = string(data)
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Println("read file error:", err)
			return
		}
		*content = string(data)
	}
}
