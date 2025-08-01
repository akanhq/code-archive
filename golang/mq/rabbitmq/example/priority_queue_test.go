package example

import (
	"fmt"
	"log"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestPriorityQueue(T *testing.T) {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	go priorityProducer(conn)
	time.Sleep(6 * time.Second)

	// 启动两个消费者
	go consumerProducer(conn, "消费者1")
	go consumerProducer(conn, "消费者2")

	select {}
}

func priorityProducer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	args := make(amqp.Table)
	args["x-max-priority"] = int32(10)
	priorityQueue, err := ch.QueueDeclare(
		"priority_queue",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatalf("声明优先级队列 失败: %v", err)
	}

	for i := 0; i < 20; i++ { // 发送 20 条消息
		msg := fmt.Sprintf("这里是消息 task_%d", i)
		err = ch.Publish(
			"",
			priorityQueue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
				Priority:    uint8(i % 10), // 优先级 0-9 循环
			},
		)
		fmt.Println("生产者发送成功", i)
		time.Sleep(100 * time.Millisecond)
	}
}

func consumerProducer(conn *amqp.Connection, consumerName string) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	args := make(amqp.Table)
	args["x-max-priority"] = int32(10)
	priorityQueue, err := ch.QueueDeclare(
		"priority_queue",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatalf("声明优先级队列 失败: %v", err)
	}

	// 设置 QoS
	err = ch.Qos(10, 0, false) // 每个消费者预取 10 条
	if err != nil {
		log.Fatalf("设置 QoS 失败: %v", err)
	}

	msgs, err := ch.Consume(
		priorityQueue.Name,
		"",
		false, // 手动确认
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("注册消费者 失败: %v", err)
	}

	for d := range msgs {
		fmt.Printf("%s 消费到消息: %s\n", consumerName, string(d.Body))
		time.Sleep(500 * time.Millisecond) // 模拟处理时间
		d.Ack(false)                       // 手动确认
	}
}
