package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

func main() {
	// Kafka 连接配置
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0 // 设置 Kafka 版本

	// 创建 Kafka 生产者
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal("Failed to start producer:", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal("Failed to close producer:", err)
		}
	}()

	// 发送消息到 Kafka
	message := &sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}
	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatal("Failed to send message:", err)
	}
	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	// 创建 Kafka 消费者
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal("Failed to close consumer:", err)
		}
	}()

	// 订阅主题
	topic := "test"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal("Failed to start partition consumer:", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal("Failed to close partition consumer:", err)
		}
	}()

	// 接收消息
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Printf("Received message: %s\n", string(msg.Value))
			case <-signals:
				fmt.Println("Interrupted")
				close(doneCh)
				return
			}
		}
	}()

	<-doneCh
	fmt.Println("Exiting...")
}
