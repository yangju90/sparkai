package main

import (
	"net/http"
	"sparkai/common/gsd"
	"sparkai/internal/handler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user/question", handler.HandleHttpRequest).Methods("POST")
	router.HandleFunc("/ws/answer", handler.HandleWebSocketConnection)

	var errChan chan (error)
	var server http.Server
	go func() {
		server := http.Server{Addr: ":8080", Handler: router}
		err := server.ListenAndServe()
		errChan <- err
	}()

	gsd.GracefulShutdown("sparkai", errChan, &server)

}
