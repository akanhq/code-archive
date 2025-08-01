package consumer_characteristic

import (
	"code_test/logger"
	"code_test/utils"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

var kafkaBrokers = []string{"localhost:9092", "localhost:9093", "localhost:9094"} // Kafka 集群地址

// Kafka 配置
const (
	replication_topic = "consumer_group_topic" // Topic 名称
	replicas          = 3                      // 复制因子数量（需要提前创建Topic时设置）
	numPartitions     = 9                      // 分区数（需要提前创建Topic时设置）
	consumersPerGroup = 3                      // 每个组的消费者数量，等于每个组负责的分区数
	deadLetterTopic   = "dead-letter-topic"
)

func TestConsumerBasic(T *testing.T) {

	logger.InitLogger("debug")

	createTopicIfNotExists(kafkaBrokers, replication_topic, replicas, numPartitions)
	createTopicIfNotExists(kafkaBrokers, deadLetterTopic, 1, 1)

	//生产者
	//go basicProducer()

	//消费者
	//go basicConsumer()

	//消费者组
	//go basicConsumerGroup("group1")
	//go basicConsumerGroup("group2")
	//go basicConsumerGroup("group3")

	//死信消费
	go deadLetterConsumerGroup()

	select {}
}

func basicConsumerGroup(groupID string) {
	group, err := createConsumerGroup(kafkaBrokers, groupID)
	if err != nil {
		logger.Errorf("创建消费者组 %s 失败: %v", groupID, err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < consumersPerGroup; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			producer, err := createSyncProducer(kafkaBrokers)
			handler := &consumerGroupHandler{
				id:       id,
				ready:    make(chan bool),
				producer: producer,
			}
			for {
				if err = group.Consume(ctx, []string{replication_topic}, handler); err != nil {
					logger.Errorf("消费者组 %s 消费失败: %v", groupID, err)
					time.Sleep(5 * time.Second)
				}
			}
		}(i)
	}
	wg.Wait()
}

func basicProducer() {
	syncProducer, err := createSyncProducer(kafkaBrokers)
	if err != nil {
		logger.Errorf("创建生产者失败 %v", err)
		return
	}
	defer syncProducer.Close()

	for i := 0; i < 100; i++ {
		id := utils.GenSnowflakeId()
		task := Task{
			ID:        id,
			Priority:  rand.Intn(10),
			Data:      fmt.Sprintf("task_%d,sort_%d", id, i+1),
			ExecuteAt: time.Now().Unix() + int64(i+1),
			Status:    "pending",
		}

		taskByte, err := json.Marshal(task)
		if err != nil {
			logger.Errorf("json序列化失败%v", err)
			continue
		}
		msg := &sarama.ProducerMessage{
			Topic: replication_topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("%c", 'A'+i%26)),
			Value: sarama.ByteEncoder(taskByte),
		}
		partition, offset, err := syncProducer.SendMessage(msg)
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			logger.Errorf("消息发送失败%v", err)
			continue
		}
		logger.Infof("消息发送成功 区间：%d 偏移量：%d Index：%d taskID：%d", partition, offset, i, task.ID)
	}
}
func basicConsumer() {
	consumer, err := createConsumer(kafkaBrokers)
	if err != nil {
		logger.Errorf("创建消费者失败 %v", err)
		return
	}
	defer consumer.Close()

	//	获取所有的分区
	partitions, _ := consumer.Partitions(replication_topic)

	//为每个分区创建一个消费者
	var wg sync.WaitGroup
	wg.Add(len(partitions))
	for _, partition := range partitions {
		go func(w *sync.WaitGroup) {
			defer w.Done()

			pc, _ := consumer.ConsumePartition(replication_topic, partition, sarama.OffsetOldest)
			defer pc.Close()

			for {
				select {
				case msg := <-pc.Messages():
					logger.Infof("接收到消息 区间 %d 偏移量 %d value %s", msg.Partition, msg.Offset, string(msg.Value))
				}
			}
		}(&wg)
	}

	wg.Wait()
}

func deadLetterConsumerGroup() {
	group, err := createConsumerGroup(kafkaBrokers, "dead-letter-topic")
	if err != nil {
		logger.Errorf("创建消费者组 %s 失败: %v", "dead-letter-topic", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			producer, err := createSyncProducer(kafkaBrokers)
			handler := &consumerGroupHandler{
				id:       id,
				ready:    make(chan bool),
				producer: producer,
			}
			for {
				if err = group.Consume(ctx, []string{deadLetterTopic}, handler); err != nil {
					logger.Errorf("消费者组 %s 消费失败: %v", "dead-letter-topic", err)
					time.Sleep(5 * time.Second)
				}
			}
		}(i)
	}
	wg.Wait()
}
