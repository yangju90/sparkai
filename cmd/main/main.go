package main

import (
	"log"
	"net/http"
	"os"
	"sparkai/common/gsd"
	"sparkai/internal/handler"

	"github.com/gorilla/mux"
)

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("当前执行目录：", dir)

	router := mux.NewRouter()
	staticDir := http.FileServer(http.Dir("D:/goconfig/static"))
	router.HandleFunc("/user/question", handler.HandleHttpRequest).Methods("POST")
	router.HandleFunc("/ws/answer", handler.HandleWebSocketConnection)
	router.PathPrefix("/image/").Handler(http.StripPrefix("/image/", staticDir))

	var errChan chan (error)
	var server http.Server
	go func() {
		server := http.Server{Addr: ":8090", Handler: router}
		err := server.ListenAndServe()
		errChan <- err
	}()

	gsd.HttpGracefulShutdown("sparkai", errChan, &server)

}
