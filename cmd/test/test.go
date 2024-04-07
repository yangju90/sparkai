package main

import (
	"fmt"
	"sync"
)

func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		ch <- i // 发送数据到通道
	}
	close(ch) // 关闭通道
}

func consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch { // 从通道接收数据
		fmt.Printf("消费者 %d 接收到数据: %d\n", id, num)
	}
}

func main() {
	// ch := make(chan int) // 创建一个整数类型的通道
	// var wg sync.WaitGroup

	// // 启动生产者
	// wg.Add(1)
	// go producer(ch, &wg)

	// // 启动两个消费者，它们共享同一个通道
	// for i := 0; i < 2; i++ {
	// 	wg.Add(1)
	// 	go consumer(i, ch, &wg)
	// }

	// wg.Wait() // 等待所有 goroutine 执行完成
}
