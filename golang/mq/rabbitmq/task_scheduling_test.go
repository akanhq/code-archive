package rabbitmq_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 配置常量
const (
	rabbitMQURL          = "amqp://admin:admin@localhost:5672/" //rabbitmq连接信息
	paymentQueueName     = "payment_events"                     //支付队列
	retryQueueName       = "payment_retries"                    //重试队列
	failedQueueName      = "payment_failed"                     //失败队列
	exchangeName         = "payment_exchange"                   //交换机名称
	maxRetryAttempts     = 3                                    //最大重试次数
	retryDelaySeconds    = 10                                   //重试等待时间
	rePublishMillisecond = 3000                                 //重新发送消息等待事件，指数避让
	rePublishAttempts    = 3                                    //最大重新发送消息次数
)

// PaymentEvent 支付事件结构体
type PaymentEvent struct {
	OrderID        string    `json:"order_id"`         //订单id
	Amount         float64   `json:"amount"`           //金额
	Status         string    `json:"status"`           //状态
	TimeStamp      time.Time `json:"time_stamp"`       //时间
	RetryCount     int       `json:"retry_count"`      //重试次数
	RetryReason    string    `json:"retry_reason"`     //重试原因
	RePublishCount int       `json:"re_publish_count"` //重新发送消息次数
}

// 用于存储待确认的消息，key 是 DeliveryTag，value 是消息内容
type PendingMessage struct {
	event *PaymentEvent
	body  []byte
}

var deliveryTag uint64 = 0 // 手动跟踪投递标签
var pendingMessages = make(map[uint64]PendingMessage)

func TestTaskScheduling(T *testing.T) {

	//	带取消的ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx

	// 启动 RabbitMQ 连接
	conn, err := setupRabbitMQConnection(ctx)
	if err != nil {
		log.Fatalf("连接 RabbitMQ 失败: %v", err)
	}
	defer conn.Close()

	//	等待组
	var wg sync.WaitGroup

	// 启动生产者（模拟支付网关）
	wg.Add(1)
	go func() {
		defer wg.Done()
		paymentProducer(ctx, conn)
	}()

	// 启动消费者（订单状态更新）
	wg.Add(2)
	go func() {
		defer wg.Done()
		paymentConsumer(ctx, conn, "worker-1")
	}()
	go func() {
		defer wg.Done()
		paymentConsumer(ctx, conn, "worker-2")
	}()

	wg.Wait()
	cancel()
	log.Println("程序已退出")
}

// setupRabbitMQConnection 带重试的 RabbitMQ 连接
func setupRabbitMQConnection(ctx context.Context) (*amqp.Connection, error) {
	for {
		conn, err := amqp.Dial(rabbitMQURL)
		if err == nil {
			go monitorConnection(ctx, conn) // 监控连接状态
			return conn, nil
		}
		log.Printf("连接 RabbitMQ 失败: %v，重试中...", err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

// monitorConnection 监控连接状态
func monitorConnection(ctx context.Context, conn *amqp.Connection) {
	closeChan := make(chan *amqp.Error)
	conn.NotifyClose(closeChan)
	select {
	case <-ctx.Done():
	case err := <-closeChan:
		log.Printf("RabbitMQ 连接断开: %v", err)
	}
}

// paymentConsumer 消费者：处理支付事件
func paymentConsumer(ctx context.Context, conn *amqp.Connection, workerID string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("创建 channel 失败: %v", err)
	}
	defer ch.Close()

	err = setupQueues(ch)
	if err != nil {
		log.Printf("设置队列 失败: %v", err)
	}

	err = ch.Qos(5, 0, false) // 每个消费者预取 5 条
	if err != nil {
		log.Printf("%s 设置 QoS 失败: %v", workerID, err)
		return
	}

	msgs, err := ch.Consume(
		paymentQueueName,
		workerID,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("%s 注册消费者失败: %v", workerID, err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s 退出", workerID)
			return
		case d, ok := <-msgs:
			if !ok {
				log.Printf("%s 消息通道关闭，尝试重连", workerID)
				ch.Close()
				break
			}

			var event PaymentEvent
			if err = json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("%s 解析消息失败: %v", workerID, err)
				d.Reject(false) // 丢弃无效消息
				continue
			}

			log.Printf("%s 处理订单: %s, 重试次数：%d", workerID, event.OrderID, event.RetryCount)
			if processPaymentEvent(&event) {
				//	处理成功
				d.Ack(false)
				log.Printf("%s 处理成功=============>: %s, 处理次数：%d", workerID, event.OrderID, event.RetryCount)
			} else if event.RetryCount < maxRetryAttempts {
				//	重试
				event.RetryCount++
				event.RetryReason = "支付状态同步失败"
				body, _ := json.Marshal(event)
				err = ch.Publish(
					exchangeName,
					"retry_key",
					false,
					false,
					amqp.Publishing{
						ContentType:  "application/json",
						Body:         body,
						DeliveryMode: amqp.Persistent,
					},
				)
				if err != nil {
					log.Printf("%s 发送重试消息失败: %v", workerID, err)
				}
				d.Ack(false) //	确认原消息
				log.Printf("%s 订单 %s 重试, 重试次数：%d", workerID, event.OrderID, event.RetryCount)
			} else {
				//	超过重试次数 发送到失败队列
				err = ch.Publish(
					exchangeName,
					"failed_key",
					false,
					false,
					amqp.Publishing{
						ContentType:  "application/json",
						Body:         d.Body,
						DeliveryMode: amqp.Persistent,
					},
				)

				if err != nil {
					log.Printf("%s 发送失败队列失败: %v", workerID, err)
				}
				d.Ack(false) //	确认原消息
				log.Printf("%s 订单 %s 重试超限，移至失败队列", workerID, event.OrderID)
			}
		}

	}

}

// paymentProducer 生产者：模拟支付事件
func paymentProducer(ctx context.Context, conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("创建 channel 失败: %v", err)
	}
	defer ch.Close()

	// 启用发布确认模式
	err = ch.Confirm(false)
	if err != nil {
		log.Printf("启用发布确认模式失败: %v", err)
		return
	}

	//注册监听确认
	acks := make(chan uint64, 1)
	nacks := make(chan uint64, 1)
	ch.NotifyConfirm(acks, nacks)

	err = setupQueues(ch)
	if err != nil {
		log.Printf("设置队列 失败: %v", err)
	}

	// 创建一个Snowflake实例，参数为workerID和dataCenterID
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		log.Printf("创建Snowflake实例 失败: %v", err)
		return
	}

	go func() {

		const batchSize = 5 // 批量确认的大小
		confirmedCount := 0 // 已确认的消息计数

		for {
			select {
			case <-ctx.Done():
				log.Println("生产者确认监听退出")
				return
			case tag, ok := <-acks:
				if !ok {
					fmt.Println("确认通过已关闭")
					return
				}
				confirmedCount++
				if _, ok := pendingMessages[tag]; ok {
					fmt.Println(fmt.Sprintf("消息确认成功 投递标签=%d", tag))
					delete(pendingMessages, tag)
				}

				//	批量确认
				if confirmedCount >= batchSize {
					log.Printf("批量确认 %d 条消息完成", confirmedCount)
					confirmedCount = 0
				}
			case tag, ok := <-nacks:
				if !ok {
					fmt.Println("失败通过已关闭")
					return
				}

				if val, ok := pendingMessages[tag]; ok {
					if val.event.RePublishCount < rePublishAttempts {
						//重新发送消息
						val.event.RePublishCount++
						body, _ := json.Marshal(val.event)
						val.body = body

						//指数避让
						time.Sleep(time.Duration(rePublishMillisecond*val.event.RePublishCount) * time.Millisecond)
						fmt.Println(fmt.Sprintf("消息确认失败 投递标签=%d,订单=%s,尝试重新发", tag, val.event.OrderID))
						err = ch.Publish(
							exchangeName,
							"payment_key",
							false,
							false,
							amqp.Publishing{
								ContentType:  "application/json",
								Body:         val.body,
								DeliveryMode: amqp.Persistent,
							},
						)
						if err != nil {
							log.Printf("发送支付事件失败: %v", err)
						} else {
							log.Printf("重发成功: 订单=%s", val.event.OrderID)
							deliveryTag++
							pendingMessages[deliveryTag] = val
						}
						//删除旧的记录
						delete(pendingMessages, tag)
					}
					//	超过最大次数不重新发送消息
				}
			}
		}
	}()

	for i := 0; i < 20; i++ {
		id, err := sf.NextID()
		if err != nil {
			log.Printf("获取雪花id 失败: %v", err)
			continue
		}

		event := PaymentEvent{
			OrderID:   fmt.Sprintf("ORDER_%d", id),
			Amount:    float64(100 + i*10),
			Status:    "pending",
			TimeStamp: time.Now(),
		}

		body, _ := json.Marshal(event)
		deliveryTag++ //手动递增投递标签
		err = ch.Publish(
			exchangeName,
			"payment_key",
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			log.Printf("发送支付事件失败: %v", err)
			continue
		}

		// 记录待确认消息
		pendingMessages[deliveryTag] = PendingMessage{
			event: &event,
			body:  body,
		}

		log.Printf("发送支付事件: %s,task_%d", event.OrderID, i)
		time.Sleep(1 * time.Second)
	}
}

// setupQueues 设置队列和交换机
func setupQueues(ch *amqp.Channel) error {

	//	声明交换机
	err := ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,  //是否持久化
		false, //是否自动删除
		false, //是否内部使用
		false, //是否等待服务器响应
		nil,   //额外参数
	)
	if err != nil {
		return fmt.Errorf("声明交换机失败: %v", err)
	}

	// 主队列（支付事件）
	_, err = ch.QueueDeclare(
		paymentQueueName,
		true,  //是否持久化
		false, //是否自动删除
		false, //是否独占
		false, //是否等待服务器响应
		nil,   //额外参数
	)
	if err != nil {
		return fmt.Errorf("声明主队列失败: %v", err)
	}
	err = ch.QueueBind(
		paymentQueueName, //队列名称
		"payment_key",    //路由键
		exchangeName,     //交换机名称
		false,            //是否等待服务器响应
		nil,              //额外参数
	)
	if err != nil {
		return fmt.Errorf("绑定主队列失败：: %v", err)
	}

	//重试队列：10秒延迟
	retryArgs := make(amqp.Table)
	retryArgs["x-message-ttl"] = int32(retryDelaySeconds * 1000)
	retryArgs["x-dead-letter-exchange"] = exchangeName
	retryArgs["x-dead-letter-routing-key"] = "payment_key"
	_, err = ch.QueueDeclare(
		retryQueueName,
		true,      //是否持久化
		false,     //是否自动删除
		false,     //是否独占
		false,     //是否等待服务器响应
		retryArgs, //额外参数
	)
	if err != nil {
		return fmt.Errorf("声明重试队列失败：: %v", err)
	}
	err = ch.QueueBind(
		retryQueueName, //队列名称
		"retry_key",    //路由键
		exchangeName,   //交换机名称
		false,          //是否等待服务器响应
		nil,            //额外参数
	)
	if err != nil {
		return fmt.Errorf("绑定重试失败：: %v", err)
	}

	// 失败队列
	_, err = ch.QueueDeclare(
		failedQueueName,
		true,  //是否持久化
		false, //是否自动删除
		false, //是否独占
		false, //是否等待服务器响应
		nil,   //额外参数
	)
	if err != nil {
		return fmt.Errorf("声明失败队列失败: %v", err)
	}
	err = ch.QueueBind(
		failedQueueName, //队列名称
		"failed_key",    //路由键
		exchangeName,    //交换机名称
		false,           //是否等待服务器响应
		nil,             //额外参数
	)
	if err != nil {
		return fmt.Errorf("绑定失败队列失败：: %v", err)
	}
	return nil
}

// processPaymentEvent 30%的失败几率
func processPaymentEvent(event *PaymentEvent) bool {
	return rand.Float64() > 0.3
}
