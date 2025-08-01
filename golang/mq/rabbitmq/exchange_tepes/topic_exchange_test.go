package exchange_types

import (
	"context"
	"fmt"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Topic Exchange
// 特点：根据路由键的模式匹配（使用通配符 * 和 #），将消息路由到符合模式的队列。
// 使用场景：需要基于模式分发消息，例如日志分级（info、error）。
func TestTopicExchangeTest(T *testing.T) {
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

	//	声明topic Exchange
	err = ch.ExchangeDeclare(
		"topic_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go ProducerTopic(ch)

	//队列1匹配所有error日志
	go ConsumerTopic(conn, "error_queue", "*.error.*", "Error Consumer")

	//队列2匹配所有us日志
	go ConsumerTopic(conn, "error_us", "#.us", "US Consumer")

	select {}
}

func ConsumerTopic(conn *amqp.Connection, queueName, match, consumerName string) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	err = ch.QueueBind(q.Name, match, "topic_exchange", false, nil)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	//	建立消费者
	msgs, err := ch.ConsumeWithContext(
		ctx,
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for d := range msgs {
		fmt.Println(fmt.Sprintf("消费者名称：%s 队列：%s 获取消息：%s", consumerName, queueName, string(d.Body)))
	}

	fmt.Println(fmt.Sprintf("%s 消费消息结束", consumerName))
}

func ProducerTopic(ch *amqp.Channel) {
	ctx := context.Background()
	var err error
	for i := 0; i < 10; i++ {

		msg, key := "", ""
		if i%2 == 0 {
			msg = fmt.Sprintf("这里是error信息 task_%d", i)
			key = "error"
		} else {
			msg = fmt.Sprintf("这里是us信息 task_%d", i)
			key = "us"
		}

		err = ch.PublishWithContext(
			ctx,
			"topic_exchange",
			key, //路由键
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("消息发送失败:", msg)
		} else {
			fmt.Println("消息发送成功:", msg)
		}
	}
}
