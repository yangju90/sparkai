package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sparkai/model/qwen"
)

func main() {
	url := "http://192.168.8.232:31434/api/chat"

	body := qwen.NewOllamaReqBody()
	bytesBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewReader(bytesBody))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	te := resp.TransferEncoding
	if te != nil && "chunked" == te[0] {
		fmt.Println("chunked!")
		reader := bufio.NewReader(resp.Body)
		for {
			chunked, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Error:", err)
				break
			}

			var respBody qwen.OllamaRespBody
			err = json.Unmarshal([]byte(chunked), &respBody)
			log.Println(chunked)
		}
	} else {
		respContent, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
		}

		var respBody qwen.OllamaRespBody
		err = json.Unmarshal([]byte(respContent), &respBody)

		log.Println(respContent)
	}

}
