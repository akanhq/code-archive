package exchange_types

import (
	"context"
	"fmt"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Headers Exchange
// 特点：根据消息头（Headers）的键值对匹配队列，而非路由键。
// 使用场景：需要基于复杂条件分发的场景，例如根据消息元数据路由。
func TestHeadersExchange(T *testing.T) {
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

	//声明Headres Exchange 交换机
	err = ch.ExchangeDeclare(
		"headers_exchange",
		"headers",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go ProducerHeaders(ch)
	go ConsumerHeaders(conn, "consumer_1", "error")
	go ConsumerHeaders(conn, "consumer_2", "info")

	select {}

}

func ConsumerHeaders(conn *amqp.Connection, consumerName, typeName string) {
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

	args := amqp.Table{
		"type":    typeName,
		"x-match": "all", // "all" 表示必须全部匹配，"any" 表示任意匹配
	}
	err = ch.QueueBind(q.Name, "", "headers_exchange", false, args)
	if err != nil {
		fmt.Println("ch.QueueBind error:", err)
		return
	}

	msgs, err := ch.ConsumeWithContext(
		context.Background(),
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("ch.ConsumeWithContext error:", err)
		return
	}

	for d := range msgs {
		fmt.Println(fmt.Sprintf("%s 消费到消息 %s", consumerName, string(d.Body)))
		time.Sleep(1 * time.Second)
	}
}

func ProducerHeaders(ch *amqp.Channel) {
	var err error
	var ctx = context.Background()

	headers := amqp.Table{
		"type":  "error",
		"level": "high",
	}
	for i := 0; i < 10; i++ {

		msg := ""

		if i%2 == 0 {
			headers["type"] = "info"
			msg = fmt.Sprintf("hello world headers exchange info test_%d", i)
		} else {
			msg = fmt.Sprintf("hello world headers exchange error test_%d", i)
			headers["type"] = "error"
		}

		err = ch.PublishWithContext(
			ctx,
			"headers_exchange",
			"", //路由键，忽略
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
				Headers:     headers,
			},
		)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("发送消息失败：", msg)
		}
		fmt.Println("消息发送成功：", msg)
	}
}
