package rabbitmq_queue

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	// 开始时间戳（毫秒），这里设置为2020-01-01 00:00:00的毫秒时间戳
	epoch int64 = 1577808000000

	// 机器ID所占的位数
	workerIDBits int64 = 5

	// 数据中心ID所占的位数
	dataCenterIDBits int64 = 5

	// 序列号所占的位数
	sequenceBits int64 = 12

	// 机器ID向左移12位
	workerIDShift int64 = sequenceBits

	// 数据中心ID向左移17位(12+5)
	dataCenterIDShift int64 = sequenceBits + workerIDBits

	// 时间戳向左移22位(5+5+12)
	timestampShift int64 = sequenceBits + workerIDBits + dataCenterIDBits

	// 生成序列号的掩码，这里为4095 (0b111111111111=0xfff=4095)
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)

	// 机器ID最大值
	maxWorkerID int64 = -1 ^ (-1 << workerIDBits)

	// 数据中心ID最大值
	maxDataCenterID int64 = -1 ^ (-1 << dataCenterIDBits)
)

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	dataCenterID  int64
	sequence      int64
}

// NewSnowflake 创建一个Snowflake实例
func NewSnowflake(workerID, dataCenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID out of range")
	}
	if dataCenterID < 0 || dataCenterID > maxDataCenterID {
		return nil, errors.New("data center ID out of range")
	}
	return &Snowflake{
		lastTimestamp: 0,
		workerID:      workerID,
		dataCenterID:  dataCenterID,
		sequence:      0,
	}, nil
}

// NextID 生成下一个唯一ID
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6 // 转换为毫秒

	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	if s.lastTimestamp == timestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 当前毫秒内的序列号已用完，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// 通过位运算生成ID
	id := ((timestamp - epoch) << timestampShift) |
		(s.dataCenterID << dataCenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id, nil
}

func TestSnowflake(T *testing.T) {
	// 创建一个Snowflake实例，参数为workerID和dataCenterID
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 生成10个ID
	for i := 0; i < 10; i++ {
		id, err := sf.NextID()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(id)
	}
}
