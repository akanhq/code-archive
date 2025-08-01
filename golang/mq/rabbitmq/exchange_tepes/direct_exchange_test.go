package exchange_types

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Direct Exchange
// 特点：根据路由键（Routing Key）精确匹配，将消息路由到绑定了相同路由键的队列。
// 使用场景：需要将消息精确发送到特定队列，例如任务分配。
func TestDirectExchange(T *testing.T) {
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

	// 声明 Direct Exchange
	err = ch.ExchangeDeclare(
		"direct_exchange", //交换机名称
		"direct",          //交换机类型
		true,              //持久化
		false,             //自动删除
		false,             //内部
		false,             //不等待
		nil,               //参数
	)
	if err != nil {
		panic(err)
	}

	go producer(ch)

	time.Sleep(3 * time.Second)
	go consumer(ch)

	select {}
}

// 生产者
func producer(ch *amqp.Channel) {

	var err error
	ctx := context.Background()
	for i := 0; i < 10; i++ {

		msg := fmt.Sprintf("Hello Direct Exchange index_%d", i)
		err = ch.PublishWithContext(
			ctx,
			"direct_exchange",
			"direct_key",
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(msg),
				DeliveryMode: amqp.Persistent, // 持久化消息
			},
		)

		time.Sleep(1 * time.Second)
		if err != nil {
			log.Println("将消息发送到交换机失败：", err)
			continue
		}
		log.Println("将消息发送到交换机成功", msg)
	}
}

// 消费者
func consumer(ch *amqp.Channel) {
	//	声明队列
	q, err := ch.QueueDeclare(
		"direct_exchange_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	//	绑定队列到交换机
	err = ch.QueueBind(
		q.Name,
		"direct_key",
		"direct_exchange",
		false, //不等待
		nil,
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	//获取消息
	msgs, err := ch.ConsumeWithContext(
		ctx,
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	for d := range msgs {
		fmt.Println("接收到消息：", string(d.Body))
		time.Sleep(1 * time.Second)
		d.Ack(true)
	}
}
