package exchange_types

import (
	"context"
	"fmt"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Fanout Exchange
// 特点：忽略路由键，将消息广播到所有绑定到该交换机的队列。
// 使用场景：需要将消息分发给多个消费者，例如日志广播、通知系统。
func TestFanoutExchange(T *testing.T) {
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

	//声明 Fanout Exchange
	err = ch.ExchangeDeclare(
		"fanout_exchange",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go ProducerFanout(ch)
	go ConsumerFanout("fanout_queue1", conn)
	go ConsumerFanout("fanout_queue2", conn)
	select {}
}

func ConsumerFanout(consumerName string, conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	err = ch.QueueBind(
		q.Name,
		"",
		"fanout_exchange",
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	msgs, err := ch.ConsumeWithContext(
		ctx,
		consumerName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(consumerName, "创建消费者失败", err)
	}

	for d := range msgs {
		fmt.Println(consumerName, "消费到消息", string(d.Body))
		time.Sleep(1 * time.Second)
	}
}

func ProducerFanout(ch *amqp.Channel) {
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Hello Fanout Exchange index+%d", i)
		err := ch.PublishWithContext(
			ctx,
			"fanout_exchange",
			"", //路由键，忽略
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)

		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("Failed to publish message", err)
		}
		fmt.Println("success to publish message", msg)
	}

}
