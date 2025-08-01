package example

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"testing"
	"time"
)

// 死信队列处理无法正常消费的消息（如被拒绝或过期），通过配置x-dead-letter-exchange实现。
func TestDeadLetterQueue(T *testing.T) {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go deadLetterQueueProducer(conn, &wg)

	time.Sleep(1 * time.Second)

	go deadLetterQueueConsumer(conn, &wg)

	wg.Wait()
}

func deadLetterQueueProducer(conn *amqp.Connection, wg *sync.WaitGroup) {

	defer wg.Done()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	//	声明死信交换机
	err = ch.ExchangeDeclare(
		"dlx_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明死信交换机: %v", err)
	}

	//	声明死信队列
	dlxQueue, err := ch.QueueDeclare(
		"dlx_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明死信队列失败: %v", err)
	}

	//	绑定死信队列到交换机
	err = ch.QueueBind(
		dlxQueue.Name,
		"dlx_key",
		"dlx_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定死信队列到死信交换机失败: %v", err)
	}

	//	声明主队列，配置死信交换机
	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = "dlx_exchange"
	args["x-dead-letter-routing-key"] = "dlx_key" // 添加 routing key
	mainQueue, err := ch.QueueDeclare(
		"main_queue",
		true,
		false,
		false,
		false,
		args,
	)

	if err != nil {
		log.Fatalf("声明主列失败: %v", err)
	}

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("这条消息是发送到主队列的 task_%d", i)
		err = ch.Publish(
			"", //默认交换机
			mainQueue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)
		time.Sleep(1 * time.Second)
		if err != nil {
			fmt.Println("发送消息到主队列失败task_", i)
			continue
		}
		fmt.Println("发送消息到主队列成功task_", i)
	}

}
func deadLetterQueueConsumer(conn *amqp.Connection, wg *sync.WaitGroup) {

	fmt.Println("开始消费")

	defer wg.Done()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// 设置 QoS，确保消费者能接收多条消息
	err = ch.Qos(10, 0, false) // prefetch count = 10
	if err != nil {
		log.Fatalf("设置 QoS 失败: %v", err)
	}

	//	声明主队列
	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = "dlx_exchange"
	args["x-dead-letter-routing-key"] = "dlx_key"
	mainQueue, err := ch.QueueDeclare(
		"main_queue",
		true,
		false,
		false,
		false,
		args, //死信参数
	)
	if err != nil {
		log.Fatalf("声明主列失败: %v", err)
	}

	//	声明死信队列
	dlxQueue, err := ch.QueueDeclare(
		"dlx_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明死信队列失败: %v", err)
	}

	//	消费主队列消息（拒绝消息）
	mainConsume, err := ch.Consume(
		mainQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("注册主消费者失败: %v", err)
	}

	go func() {
		for d := range mainConsume {
			fmt.Println("主队消费到消息", string(d.Body))
			d.Reject(false) // 拒绝消息，进入死信队列
			time.Sleep(500 * time.Millisecond)
		}
	}()

	dlxConsume, err := ch.Consume(
		dlxQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("注册死信消费者失败: %v", err)
	}

	go func() {
		for d := range dlxConsume {
			fmt.Println("死信消费到消息", string(d.Body))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	time.Sleep(15 * time.Second)
}
