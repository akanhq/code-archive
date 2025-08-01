package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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

	//	声明交换机
	err = ch.ExchangeDeclare(
		"direct_exchange", //交换机名称
		"direct",          //交换机类型
		true,              //是否持久化
		false,             //是否自动删除
		false,             //是否内部,
		false,             //是否阻塞
		nil,               //额外参数
	)

	//	声明队列
	q, err := ch.QueueDeclare(
		"direct_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败: %v", err)
	}

	// 绑定队列到交换机
	err = ch.QueueBind(
		q.Name,            // 队列名称
		"direct_routing",  // 路由键
		"direct_exchange", // 交换机名称
		false,             // 是否阻塞
		nil,               // 额外参数
	)
	if err != nil {
		log.Fatalf("绑定队列失败: %v", err)
	}
	// 消费消息
	ctx := context.Background()
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

	// 接收并处理消息
	for d := range msgs {
		log.Printf("收到消息: %s", d.Body)
	}
}
