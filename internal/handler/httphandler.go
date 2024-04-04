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
	var resp *model.HttpBodyResponse

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		str, _ := json.Marshal(body)
		if v, ok := mem.WSConnContainers[userId]; ok {
			v.Send(str)
			resp = model.Success()
		} else {
			resp = model.UserIdNotOnline(userId)
		}
		res, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(res))
	} else {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	}
}
