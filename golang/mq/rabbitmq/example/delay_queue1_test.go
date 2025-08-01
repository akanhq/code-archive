package example

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func TestDelayQueue1(T *testing.T) {

	// 连接到 RabbitMQ
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//	声明延迟交换机
	args := amqp.Table{
		"x-delayed-type": "direct", //指定延迟交换机的类型
	}
	err = ch.ExchangeDeclare(
		"delayed_ecchange",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare an exchange")

	//	声明队列
	q, err := ch.QueueDeclare(
		"delayed_queue", //队列名称
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	//	绑定队列到交换机
	err = ch.QueueBind(q.Name, "delay_key", "delayed_ecchange", false, nil)
	failOnError(err, "Failed to bind a queue")

	//	生产者，发送延迟消息
	go func() {

		chProducer, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer chProducer.Close()

		headers := amqp.Table{"x-delay": int32(5000)} // 延迟5秒

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("hello this is a delayed message task_%d", i)
			err = chProducer.PublishWithContext(
				ctx,
				"delayed_ecchange",
				"delay_key",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg),
					Headers:     headers,
				},
			)

			if err != nil {
				fmt.Println("消息发送失败", i)
				continue
			}
			fmt.Println("消息发送成功", time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	//	消费者，接收消息
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	for d := range msgs {
		fmt.Println("接收到消息", string(d.Body), time.Now().Format("2006-01-02 15:04:05"))
	}

}
