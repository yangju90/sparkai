package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sparkai/internal/io"
	"sparkai/model"
	"sparkai/model/constant"
	"sparkai/model/mem"
	"strings"
)

func HandleHttpRequest(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	userId := query.Get("userId")

	var body model.HttpBodyRequest
	var resp *model.HttpBodyResponse

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		if v, ok := mem.WSConnContainers[userId]; ok {
			if len(body.Text) != 0 {
				// 1. 添加用户提问
				v.AppendMessage(body.Text, constant.USER)
				v.ImageData = body.ImageData

				// 1.调用sparkai
				if e := io.Wsservice(userId); e != nil {
					resp = model.Faild(e)
				} else {
					resp = model.Success()
				}
			}
		} else {
			resp = model.UserIdNotOnline(userId)
		}
		res, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(res))
	} else {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	}
}
