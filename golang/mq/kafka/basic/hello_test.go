package basic

import (
	"code_test/logger"
	"fmt"
	"github.com/IBM/sarama"
	"testing"
	"time"
)

func TestHelloKafka(T *testing.T) {

	// 配置生产者
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // 确保消息发送成功后返回确认

	logger.InitLogger("debug")
	//go Producer(config)
	go Consumer(config)
	select {}
}

func Consumer(config *sarama.Config) {
	//	创建消费者
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		logger.Errorf("创建消费者失败")
	}
	defer consumer.Close()

	//	订阅topic
	partitionConsumer, err := consumer.ConsumePartition("test-topic", 0, sarama.OffsetOldest)
	if err != nil {
		logger.Errorf("订阅topic失败")
	}
	defer partitionConsumer.Close()

	//	消费消息
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			logger.Debugf("消费到消息%s", string(msg.Value))
		case err = <-partitionConsumer.Errors():
			logger.Errorf("消费消息失败：%s", err.Error())
		}
		time.Sleep(1 * time.Second)
	}

}

func Producer(config *sarama.Config) {

	// 创建生产者
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		logger.Errorf("创建生产者失败")
	}
	defer producer.Close()

	for i := 0; i < 10; i++ {
		//	发送消息
		msg := &sarama.ProducerMessage{
			Topic: "test-topic",
			Value: sarama.StringEncoder(fmt.Sprintf("hello world_%d", i)),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			logger.Errorf("发送消息失败")
		}

		logger.Infof("消息发送成功 partition %d at offset %d\n", partition, offset)
		time.Sleep(1 * time.Second)
	}
}
