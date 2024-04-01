package gsd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown(appName string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Main goroutine is waiting for signal...")
	// 等待信号
	sig := <-c
	log.Printf("Received signal: %v\n", sig)

	log.Println("Cleaning up...")
	// 这里可以执行一些清理工作，比如关闭文件、释放资源等

	log.Printf("%s program exit ok\n", appName)
}
