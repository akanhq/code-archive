package example

import (
	"fmt"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 消息确认机制确保消息被消费者成功处理。通过手动确认（Manual Acknowledgment），消费者可以在处理完消息后再通知RabbitMQ，避免消息丢失。
func TestPublisherConfirmsConsumerAck(t *testing.T) {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		t.Fatalf("无法连接到 RabbitMQ: %v", err)
	}
	defer conn.Close()

	go publisherConfirmsProducer(conn)
	go publisherConfirmsConsumer(conn)

	select {}
}

func publisherConfirmsProducer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	//	声明一个队列
	q, err := ch.QueueDeclare(
		"ack_queue",
		false, // 队列名称
		false, // 非持久化
		false, // 不自动删除
		false, // 非排他性
		nil,   // 不等待
	)
	if err != nil {
		panic(err)
	}

	//注册监听确认
	acks := make(chan uint64, 1)
	nacks := make(chan uint64, 1)
	ch.NotifyConfirm(acks, nacks)

	//监听消息是否发送到队列
	go func() {
		for {
			select {
			case tag, ok := <-acks:
				if !ok {
					fmt.Println("确认通过已关闭")
					return
				}
				fmt.Println(fmt.Sprintf("消息确认成功：投递标签=%d", tag))
			case tag, ok := <-nacks:
				if !ok {
					fmt.Println("失败通过已关闭")
					return
				}
				fmt.Println(fmt.Sprintf("消息确认失败：投递标签=%d", tag))
			}
		}
	}()

	//	往队列发送消息
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("this is msg_%d", i)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("消息发送失败 ==========> ", msg)
			continue
		}
		fmt.Println("消息发送成功 ==========> ", msg)
	}
}

func publisherConfirmsConsumer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	//	声明一个队列
	q, err := ch.QueueDeclare(
		"ack_queue",
		false, // 队列名称
		false, // 非持久化
		false, // 不自动删除
		false, // 非排他性
		nil,   // 不等待
	)
	if err != nil {
		panic(err)
	}

	//	创建消费者
	msgs, err := ch.Consume(
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

	for msg := range msgs {
		fmt.Println("消费到消息 ==========> ", string(msg.Body))
		msg.Ack(false)
		time.Sleep(1 * time.Second)
	}
}
