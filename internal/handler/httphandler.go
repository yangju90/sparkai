package handler

import (
	"fmt"
	"net/http"
	"strings"
)

func HandleHttpRequest(w http.ResponseWriter, r *http.Request) {
	// 处理POST请求
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		// 解析JSON数据并进行处理
		// ...
		fmt.Fprintf(w, "收到POST请求")
	} else {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	}
}
