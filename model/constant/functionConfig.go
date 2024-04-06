package constant

import (
	"embed"
	"encoding/json"
	"log"
	"os"
	resources "sparkai/resources"
)

var FunctionsConfig map[string]interface{}

func init() {
	log.Println("init FunctionsConfig ......")

	path := "E:/goconfig/function" 
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Println("Read local file error:", err)
		return
	}

	if len(entries) == 0 {
		log.Println("read default files...")
		FunctionsConfig, err = ReadFunctionsConfig(&resources.FunctionsResource)
	} else {
		log.Println("read local files...")
		FunctionsConfig, err = ReadLocalFunctionsConfig(entries, path)
	}

	if err != nil {
		log.Println("init FunctionsConfig error!")
	}

}

func ReadLocalFunctionsConfig(entries []os.DirEntry, path string) (configMsg map[string]interface{}, err error) {

	var regfuncs []map[string]interface{}

	for _, v := range entries {
		var data []byte
		data, err = os.ReadFile(path + "/" + v.Name())
		if err != nil {
			log.Println("read local file error:", v.Name(), err)
			return
		}
		var funcMsg map[string]interface{}

		err = json.Unmarshal(data, &funcMsg)
		if err != nil {
			log.Println("Error parsing config file:", err)
			return
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

func ReadFunctionsConfig(resources *embed.FS) (configMsg map[string]interface{}, err error) {
	path := "function"
	dir, err := resources.ReadDir(path)
	if err != nil {
		log.Println("Error reading config file:", err)
		return
	}

	var regfuncs []map[string]interface{}

	for _, v := range dir {
		vPath := path + "/" + v.Name()
		data, err := resources.ReadFile(vPath)
		if err != nil {
			log.Println("Error reading config file:", err)
			break
		}
		var funcMsg map[string]interface{}

		err = json.Unmarshal(data, &funcMsg)
		if err != nil {
			log.Println("Error parsing config file:", err)
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
