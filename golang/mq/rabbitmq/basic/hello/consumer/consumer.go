package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	// 连接到 RabbitMQ 服务器
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("打开通道失败: %v", err)
	}
	defer ch.Close()

	// 声明一个队列
	q, err := ch.QueueDeclare(
		"hello", // 队列名称
		true,    // 是否持久化
		false,   // 是否自动删除
		false,   // 是否独占
		false,   // 是否阻塞
		nil,     // 额外参数
	)
	if err != nil {
		log.Fatalf("声明队列失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("开始时间", time.Now().Format("2006-01-02 15:04:05"))

	// 消费消息
	msgs, err := ch.ConsumeWithContext(
		ctx,
		q.Name, // 队列名称
		"",     // 消费者标签
		true,   // 自动确认
		false,  // 是否独占
		false,  // 是否本地
		false,  // 是否阻塞
		nil,    // 额外参数
	)
	if err != nil {
		log.Fatalf("注册消费者失败: %v", err)
	}

	for d := range msgs {
		fmt.Println("消费到消息：", string(d.Body), "时间：", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(1 * time.Second)
	}
	fmt.Println("消费结束", time.Now().Format("2006-01-02 15:04:05"))
}
