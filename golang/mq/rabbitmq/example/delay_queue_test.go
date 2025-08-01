package example

import (
	"fmt"
	"log"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestDelayQueue(T *testing.T) {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	go delayQueueProducer(conn)
	go delayQueueConsumer(conn)

	select {}
}

func delayQueueProducer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// 声明交换机
	err = ch.ExchangeDeclare(
		"delay_exchange",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	// 声明延迟队列，设置TTL和死信交换机
	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = "delay_exchange"
	args["x-dead-letter-routing-key"] = "target_routing_key"
	args["x-message-ttl"] = int32(5000) // 5秒延迟
	_, err = ch.QueueDeclare(
		"delay_queue",
		true,
		false,
		false,
		false,
		args, // 修正：传入 args
	)
	if err != nil {
		log.Fatalf("声明延迟队列 失败: %v", err)
	}

	// 声明目标队列
	_, err = ch.QueueDeclare(
		"target_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明目标队列 失败: %v", err)
	}

	// 绑定延迟队列到交换机，使用空路由键匹配生产者
	err = ch.QueueBind(
		"delay_queue",
		"", // 与生产者的路由键一致
		"delay_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定延迟队列到交换机 失败: %v", err)
	}

	// 绑定目标队列到交换机，使用死信路由键
	err = ch.QueueBind(
		"target_queue",
		"target_routing_key", // 与死信路由键一致
		"delay_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定目标队列到交换机 失败: %v", err)
	}

	// 发送延迟消息
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("这是延迟消息 task_%d", i)
		err = ch.Publish(
			"delay_exchange",
			"", // 空路由键，路由到 delay_queue
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		if err != nil {
			fmt.Println("发送消息到交换机 失败", err, "msg=", i)
			continue
		}
		fmt.Println("发送消息成功", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

func delayQueueConsumer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// 声明目标队列
	targetQueue, err := ch.QueueDeclare(
		"target_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明目标队列 失败: %v", err)
	}

	// 注册目标消费者
	targetConsume, err := ch.Consume(
		targetQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("注册目标消费者 失败: %v", err)
	}
	go func() {
		for msg := range targetConsume {
			fmt.Println("目标队列收到消息：", string(msg.Body), time.Now().Format("2006-01-02 15:04:05"))
		}
	}()
	select {}
}
