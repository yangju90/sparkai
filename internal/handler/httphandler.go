package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sparkai/model"
	"sparkai/model/mem"
	"strings"
)

func HandleHttpRequest(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	userId := query.Get("userId")

	var body model.HttpBodyRequest

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		// 解析JSON数据并进行处理
		// ...
		str, _ := json.Marshal(body)
		if v, ok := mem.WSConnContainers[userId]; ok {

			v.Send(str)
		}
		fmt.Fprintf(w, "收到POST请求"+string(str))
	} else {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	}
}
