package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	msg := "hello RabbitMQ"
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 模拟延迟，确保发送时上下文已过期
	time.Sleep(3 * time.Second)

	for {
		err = ch.PublishWithContext(
			ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		if err != nil {
			fmt.Println("发送消息失败：", err)
		} else {
			fmt.Println("发送消息成功")
		}
		time.Sleep(2 * time.Second)
	}
}
