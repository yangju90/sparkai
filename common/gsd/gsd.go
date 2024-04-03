package gsd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(appName string, errChan chan (error), server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Main goroutine is waiting for signal...")
	// 等待信号
	var err error
	var sig os.Signal
	select {
	case err = <-errChan:
		log.Printf("Received starting web signal: %v\n", err)
	case sig = <-c:
		log.Printf("Received signal: %v\n", sig)
		log.Printf("%s web service is exiting... \n", appName)
		log.Println("Cleaning up...")
		// 这里可以执行一些清理工作，比如关闭文件、释放资源等
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		server.Shutdown(ctx) // 优雅关闭http服务实例

		log.Printf("%s program exit ok\n", appName)
	}

}
