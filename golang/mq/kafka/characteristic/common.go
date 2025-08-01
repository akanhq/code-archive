package characteristic

import (
	"code_test/logger"
	"github.com/IBM/sarama"
	"sync"
)

// 在使用 Sarama 库（Kafka 的 Go 客户端）实现消费者组时，需要提供一个实现了 sarama.ConsumerGroupHandler 接口的结构体。
// 这个接口定义了消费者组在消费消息时需要实现的
type consumerGroupHandler struct {
	ready chan bool
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	//标记消费者，准备就绪
	close(h.ready)
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		logger.Infof("消费者组 %s 接收到消息 topic:%s, 分区:%d, 偏移量:%d, value:%s\n", session.MemberID(), message.Topic, message.Partition, message.Offset, string(message.Value))

		//标记消息为已处理（自动提交偏移量）
		session.MarkMessage(message, "")

		//批量提交
		if message.Offset%100 == 0 {
			session.Commit() // 显式提交偏移量
		}

		//time.Sleep(1 * time.Second)
	}
	return nil
}

// CustomKeyPartitioner 实现完全自定义的分区逻辑
type CustomKeyPartitioner struct {
	partition int32 // 用于轮询策略的状态
	mu        sync.Mutex
}

func (p *CustomKeyPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	key, _ := message.Key.Encode()
	partition := int32(key[0] % byte(numPartitions))
	//log.Printf("KEY: %s → 分区%d", string(key), partition)
	return partition, nil
}

func (p *CustomKeyPartitioner) RequiresConsistency() bool {
	return true
}

func (p *CustomKeyPartitioner) Close() error {
	return nil
}
