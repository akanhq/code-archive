package example

import (
	"fmt"
	"log"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 消息持久化确保消息在RabbitMQ服务器重启后不会丢失，需要队列和消息都设置为持久化。
func TestPersistence(T *testing.T) {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	go persistenceProducer(conn)
	go persistenceConsumer(conn)

	select {}
}

func persistenceProducer(conn *amqp.Connection) {

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	//	声明持久化队列
	q, err := ch.QueueDeclare(
		"persistence_queue",
		true, //持久化队列
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("this is persistence message %d", i)

		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(msg),
				DeliveryMode: amqp.Persistent,
			},
		)

		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println("Persistence message sent", msg)
	}

}
func persistenceConsumer(conn *amqp.Connection) {

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	//	声明持久化队列
	q, err := ch.QueueDeclare(
		"persistence_queue",
		true, //持久化队列
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		fmt.Println("消费到消息", string(msg.Body))
		time.Sleep(1 * time.Second)
		msg.Ack(false)
	}
}
