package characteristic

//import (
//	"code_test/logger"
//	"context"
//	"fmt"
//	"github.com/IBM/sarama"
//	"log"
//	"sync"
//	"testing"
//	"time"
//)
//
//var kafkaBrokers = []string{"localhost:9092", "localhost:9093", "localhost:9094"} // Kafka 集群地址
//
//// Kafka 配置
//const (
//	replication_topic = "replicated-topic" // Topic 名称
//	replicas          = 3                  // 复制因子数量（需要提前创建Topic时设置）
//)
//
//// Kafka 通过复制（Replication）确保数据可靠性。如果一个 Broker 失败，其他副本可以接管。以下是关键概念：
//// 复制因子（Replication Factor）：指定每个分区的副本数（例如 2 或 3）。副本包括领导者（Leader）和跟随者（Follower）。
//// 领导者：处理所有读写请求的副本。
//// 跟随者：同步领导者的数据，但不处理客户端请求。
//func TestReplication(T *testing.T) {
//	logger.InitLogger("debug")
//
//	//创建topic
//	createTopicIfNotExists(kafkaBrokers, replication_topic, replicas, 3)
//
//	go replicationProducer()
//	//go replicationConsumer()
//	go replicationConsumerGroup()
//
//	select {}
//}
//
//func replicationProducer() {
//	producer := startProducer()
//	defer producer.Close()
//	for i := 0; ; i++ {
//		msg := &sarama.ProducerMessage{
//			Topic: replication_topic,
//			Value: sarama.StringEncoder(fmt.Sprintf("replication message task_%d", i)),
//		}
//		partition, offset, err := producer.SendMessage(msg)
//		if err != nil {
//			logger.Errorf("消息发送失败 -> 分区: %d, 偏移量: %d err=%v", partition, offset, err)
//		} else {
//			logger.Infof("消息发送成功 -> 分区: %d, 偏移量: %d", partition, offset)
//		}
//		time.Sleep(1 * time.Second)
//	}
//}
//
//func replicationConsumerGroup() {
//	go startConsumerGroup("消费者组1")
//	go startConsumerGroup("消费者组2")
//	go startConsumerGroup("消费者组3")
//	select {}
//}
//func replicationConsumer() {
//	consumer := startConsumer()
//	defer consumer.Close()
//
//	//	获取分区列表
//	partitions, err := consumer.Partitions(replication_topic)
//	if err != nil {
//		fmt.Errorf("获取分区列表失败 -> %s", err.Error())
//	}
//
//	var wg sync.WaitGroup
//	wg.Add(len(partitions))
//
//	//	为每个分区创建消费者
//	for _, partition := range partitions {
//		go func(p int32) {
//			defer wg.Done()
//
//			pc, err := consumer.ConsumePartition(replication_topic, p, sarama.OffsetNewest)
//			if err != nil {
//				logger.Errorf("err =%v", err)
//				return
//			}
//			defer pc.Close()
//			for {
//				select {
//				case msg := <-pc.Messages():
//					logger.Infof("分区 %d 偏移量 %d 接收到消息%s", msg.Partition, msg.Offset, string(msg.Value))
//
//				}
//			}
//
//		}(partition)
//	}
//	wg.Wait()
//
//}
//
//func startProducer() sarama.SyncProducer {
//
//	config := sarama.NewConfig()
//	config.Producer.RequiredAcks = sarama.WaitForAll //等待所有副本确认
//	config.Producer.Retry.Max = 5
//	config.Producer.Return.Successes = true
//
//	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
//	if err != nil {
//		logger.Errorf("创建生产者失败 err=%v", err)
//	}
//	return producer
//}
//
//// startConsumerGroup 创建消费者组
//func startConsumerGroup(groupID string) {
//	config := sarama.NewConfig()
//	config.Consumer.Group.Session.Timeout = 10 * time.Second
//	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
//
//	consumerGroup, err := sarama.NewConsumerGroup(kafkaBrokers, groupID, config)
//	if err != nil {
//		logger.Errorf("创建消费者组失败 err=%v", err)
//		return
//	}
//	defer consumerGroup.Close()
//
//	//	消费逻辑
//	ctx := context.Background()
//	topics := []string{replication_topic}
//
//	handle := &consumerGroupHandler{ready: make(chan bool)}
//
//	for {
//		//每次循环都创建新的ready通道
//		handle.ready = make(chan bool)
//		err = consumerGroup.Consume(ctx, topics, handle)
//		if err != nil {
//			logger.Errorf("消费错误 err=%v", err)
//		}
//		if ctx.Err() != nil {
//			return
//		}
//	}
//}
//
//// 消费者组处理器
////type consumerGroupHandler struct{}
//
//// 创建消费者
//func startConsumer() sarama.Consumer {
//	config := sarama.NewConfig()
//	config.Consumer.Return.Errors = true
//
//	consumer, err := sarama.NewConsumer(kafkaBrokers, config)
//	if err != nil {
//		logger.Errorf("创建消费者失败 err=%v", err)
//	}
//	return consumer
//}
//
//func createTopicIfNotExists(brokerList []string, topic string, replicationFactor, numPartitions int32) {
//
//	config := sarama.NewConfig()
//
//	// 创建 Admin 客户端
//	admin, err := sarama.NewClusterAdmin(brokerList, config)
//	if err != nil {
//		logger.Errorf("创建admin客户端失败 err=%v", err)
//		return
//	}
//	defer admin.Close()
//
//	// 获取并打印所有broker元数据
//	if brokers, _, err := admin.DescribeCluster(); err == nil {
//		logger.Infof("集群broker信息: %+v", brokers)
//	} else {
//		logger.Errorf("获取broker信息失败: %v", err)
//	}
//
//	// 检查 Topic 是否存在
//	topics, err := admin.ListTopics()
//	if err != nil {
//		log.Fatalf("获取topics 失败 err=%v", err)
//		return
//	}
//
//	_, exists := topics[topic]
//	if !exists {
//
//		// 定义 Topic 详情
//		detail := &sarama.TopicDetail{
//			NumPartitions:     numPartitions,            // 分区数
//			ReplicationFactor: int16(replicationFactor), // 复制因子
//			ConfigEntries:     map[string]*string{},     // 可选配置
//		}
//
//		// 创建 Topic
//		err = admin.CreateTopic(topic, detail, false)
//		if err != nil {
//			if err.Error() == "kafka server: Topic with this name already exists" {
//				logger.Warnf("Topic %s 已经存在，跳过创建", topic)
//			} else {
//				logger.Debugf("创建topic失败 err=%v", err)
//			}
//		} else {
//			logger.Infof("Topic %s 创建成功 复制因子 %d 分区数 %d\n", topic, replicationFactor, numPartitions)
//		}
//	} else {
//		fmt.Printf("Topic %s 已经存在\n", topic)
//	}
//}
