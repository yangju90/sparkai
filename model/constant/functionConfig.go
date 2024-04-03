package constant

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	resources "sparkai/resources"
)

var FunctionsConfig map[string]interface{}

func init() {
	log.Println("init FunctionsConfig ......")

	var err error
	FunctionsConfig, err = ReadFunctionsConfig(&resources.FunctionsResource)

	if err != nil {
		log.Println("init FunctionsConfig error!")
	}
}

func ReadFunctionsConfig(resources *embed.FS) (configMsg map[string]interface{}, err error) {
	path := "function"
	dir, err := resources.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var regfuncs []map[string]interface{}

	for _, v := range dir {
		vPath := path + "/" + v.Name()
		data, err := resources.ReadFile(vPath)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			break
		}
		var funcMsg map[string]interface{}

		err = json.Unmarshal(data, &funcMsg)
		if err != nil {
			fmt.Println("Error parsing config file:", err)
			break
		}
		regfuncs = append(regfuncs, funcMsg)
	}

	if err == nil {
		configMsg = map[string]interface{}{
			"text": regfuncs,
		}
	}

	return
}
