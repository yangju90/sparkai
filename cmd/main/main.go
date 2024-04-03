package main

import (
	"log"
	"net/http"
	"sparkai/common/gsd"
	"sparkai/internal/handler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user/question", handler.HandleHttpRequest).Methods("POST")
	router.HandleFunc("/ws/answer", handler.HandleWebSocketConnection)

	log.Fatal(http.ListenAndServe(":8080", router))

	gsd.GracefulShutdown("sparkai")
}
