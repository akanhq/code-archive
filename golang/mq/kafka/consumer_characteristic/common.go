package consumer_characteristic

import (
	"code_test/logger"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// 在使用 Sarama 库（Kafka 的 Go 客户端）实现消费者组时，需要提供一个实现了 sarama.ConsumerGroupHandler 接口的结构体。
// 这个接口定义了消费者组在消费消息时需要实现的
type consumerGroupHandler struct {
	id       int
	ready    chan bool
	producer sarama.SyncProducer //死信队列
}

func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(h.ready)
	logger.Infof("消费者组 %s 准备就绪，分配的分区: %+v", session.MemberID(), session.Claims())
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	logger.Errorf("Rebalance即将发生，释放分区资源")
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		logger.Infof("消费者组 %s 接收到消息 topic:%s, 分区:%d, 偏移量:%d, value:%s\n", session.MemberID(), message.Topic, message.Partition, message.Offset, string(message.Value))

		//标记消息为已处理（自动提交偏移量）
		session.MarkMessage(message, "")
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// ConsumeClaim 处理消息：死信队列
//func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
//	for message := range claim.Messages() {
//
//		logger.Infof("消费者组 %s 接收到消息 topic:%s, 分区:%d, 偏移量:%d, value:%s\n", session.MemberID(), message.Topic, message.Partition, message.Offset, string(message.Value))
//		if err := h.processMessage(message); err != nil {
//			//	处理失败，将消息发送到死信队列
//			if err = h.sendToDeadLetter(message); err != nil {
//				logger.Errorf("%v", err)
//			}
//			continue
//		}
//		//标记消息为已处理（自动提交偏移量）
//		session.MarkMessage(message, "")
//	}
//
//	return nil
//}

func (h *consumerGroupHandler) sendToDeadLetter(message *sarama.ConsumerMessage) error {

	_, _, err := h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "dead-letter-topic",
		Value: sarama.ByteEncoder(message.Value),
		Key:   sarama.ByteEncoder(message.Key),
	})
	if err != nil {
		return fmt.Errorf("failed to send to dead letter topic: %v", err)
	}
	logger.Infof("死信队列消息发送成功 value:%v", string(message.Value))
	return nil
}

func (h *consumerGroupHandler) processMessage(message *sarama.ConsumerMessage) error {
	return fmt.Errorf("订单异常处理失败")
}

// createSyncProducer 创建 Kafka 同步生产者
func createSyncProducer(kafkaBrokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //等待所有副本确认
	config.Producer.Retry.Max = 5                    //生产者发送消息失败时的最大重试次数
	config.Producer.Return.Successes = true          //发送消息后会阻塞等待响应，需要Successes通道来接收成功确认
	config.Net.MaxOpenRequests = 1                   //允许的未完成请求的最大数量，在同步生产者模式下，MaxOpenRequests 必须为 1，以确保消息严格按顺序发送
	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	return producer, err
}

func createConsumerGroup(kafkaBrokers []string, groupId string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest //从最早的开始消费

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRange(),
	}

	group, err := sarama.NewConsumerGroup(kafkaBrokers, groupId, config)
	return group, err
}

func createConsumer(kafkaBrokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	consumer, err := sarama.NewConsumer(kafkaBrokers, config)
	return consumer, err
}

// createTopicIfNotExists 创建topic
func createTopicIfNotExists(brokerList []string, topic string, replicationFactor, numPartitions int32) {

	config := sarama.NewConfig()

	// 创建 Admin 客户端
	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		logger.Errorf("创建admin客户端失败 err=%v", err)
		return
	}
	defer admin.Close()

	// 获取并打印所有broker元数据
	if brokers, _, err := admin.DescribeCluster(); err == nil {
		brokerStr := ""
		for _, broker := range brokers {
			brokerStr += fmt.Sprintf(",%s", broker.Addr())
		}
		logger.Infof("集群broker信息: %v", brokerStr)
	} else {
		logger.Errorf("获取broker信息失败: %v", err)
	}

	// 检查 Topic 是否存在
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("获取topics 失败 err=%v", err)
		return
	}

	_, exists := topics[topic]
	if !exists {

		// 定义 Topic 详情
		detail := &sarama.TopicDetail{
			NumPartitions:     numPartitions,            // 分区数
			ReplicationFactor: int16(replicationFactor), // 复制因子
			ConfigEntries:     map[string]*string{},     // 可选配置
		}

		// 创建 Topic
		err = admin.CreateTopic(topic, detail, false)
		if err != nil {
			if err.Error() == "kafka server: Topic with this name already exists" {
				logger.Warnf("Topic %s 已经存在，跳过创建", topic)
			} else {
				logger.Debugf("创建topic失败 err=%v", err)
			}
		} else {
			logger.Infof("Topic %s 创建成功 复制因子 %d 分区数 %d\n", topic, replicationFactor, numPartitions)
		}
	} else {
		logger.Infof("Topic %s 已经存在\n", topic)
	}
}
