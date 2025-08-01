package characteristic

import (
	"code_test/logger"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

//// 在使用 Sarama 库（Kafka 的 Go 客户端）实现消费者组时，需要提供一个实现了 sarama.ConsumerGroupHandler 接口的结构体。
//// 这个接口定义了消费者组在消费消息时需要实现的
//type consumerGroupHandler struct {
//	ready chan bool
//}
//
//func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
//	//标记消费者，准备就绪
//	close(h.ready)
//	return nil
//}
//
//func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
//	return nil
//}
//
//func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
//
//	for message := range claim.Messages() {
//		logger.Infof("消费者组 %s 接收到消息 topic:%s, 分区:%d, 偏移量:%d, value:%s\n",
//			session.MemberID(), message.Topic, message.Partition, message.Offset, string(message.Value))
//		session.MarkMessage(message, "")
//		session.Commit() // 显式提交偏移量
//	}
//
//	return nil
//}

func TestConsumerGroupTest(t *testing.T) {
	logger.InitLogger("debug")

	// 启动生产者
	go groupProducer()

	// 启动消费者组
	go groupConsumer()

	select {}
}

func groupConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //从最早的开始消费

	//创建消费者组
	group, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "test-consumer-group", config)
	if err != nil {
		logger.Errorf("创建消费者组失败 err %v", err)
	}

	defer group.Close()
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)

		handle := &consumerGroupHandler{ready: make(chan bool)}
		ctx := context.Background()
		go func(i int) {
			defer wg.Done()

			//每次循环都创建新的ready通道
			handle.ready = make(chan bool)
			err = group.Consume(ctx, []string{"test-topic"}, handle)
			if err != nil {
				logger.Errorf("Error from consumer %d: %v", i, err)
				return
			}

			// 检查上下文是否被取消
			if ctx.Err() != nil {
				return
			}

			// 等待消费者准备就绪
			<-handle.ready
		}(i)
	}
	wg.Wait()
	fmt.Println("消费结束")
}
func groupProducer() {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //等待所有副本确认
	config.Producer.Retry.Max = 5                    //最大重试次数
	config.Producer.Return.Successes = true

	//	创建同步生产者
	//sarama.NewSyncProducer，这是同步生产者。每次发送消息都会阻塞，直到 Kafka 响应。
	//这确保了消息的可靠性，但性能可能不如异步生产者（sarama.NewAsyncProducer）。
	//如果需要更高吞吐量，可以考虑使用异步生产者，但需要处理更多的错误情况
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		logger.Errorf("创建同步生产者失败 err %v", err)
	}
	defer producer.Close()

	for i := 0; i < 1000; i++ {
		//Kafka 会根据键的哈希值决定分区，确保相同键的消息总是发送到同一个分区。
		//如果消息指定了键（Key），Kafka 会使用一个哈希函数（通常是 Murmur2 哈希）来决定消息应该发送到哪个分区。分区分配的公式如下
		//partition = hash(key) mod numPartitions
		msg := &sarama.ProducerMessage{
			Topic: "test-topic",
			Key:   sarama.StringEncoder(fmt.Sprintf("key_%d", i)),
			Value: sarama.StringEncoder(fmt.Sprintf("this is msg in test-topic task_%d", i)),
		}
		partition, offset, err := producer.SendMessage(msg)
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			logger.Errorf("发送消息失败 err %v", err)
			continue
		}
		logger.Debugf("发送消息成功 分区: %d,偏移量：%d", partition, offset)
	}
}
