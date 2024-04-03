package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HandleHttpRequest(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取请求体失败：", err)
		return
	}
	fmt.Println("请求体内容：", string(body))

	// dec := json.NewDecoder(r.Body)
	// if err := dec.Decode(&book); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// 处理POST请求
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		// 解析JSON数据并进行处理
		// ...
		fmt.Fprintf(w, "收到POST请求")
	} else {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	}
}
