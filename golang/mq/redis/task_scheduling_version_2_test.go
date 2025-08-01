package main

import (
	"code_test/logger"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

const (
	TasksPending    = "tasks_pending"     //调度队列(Sorted Set)
	TaskMetadata    = "task_metadata"     //任务元数据(Hash)
	TaskQueueFailed = "task_queue_failed" //失败队列(Sorted Set)-失败次数超过指定次数
)

var workerStreams = []string{"worker_stream_1", "worker_stream_2"}

// Task 定义任务结构体
type Task struct {
	ID         uint64 `json:"id"`
	Priority   int    `json:"priority"`    // 优先级，越大越优先
	Data       string `json:"data"`        // 任务内容
	ExecuteAt  int64  `json:"execute_at"`  // 执行时间戳
	Status     string `json:"status"`      // 状态: pending, processing, success, failed
	RetryCount int    `json:"retry_count"` // 重试次数
	StartedAt  int64  `json:"started_at"`  // 开始处理时间戳（用于超时检测）
}

func TestTaskSchedulingVersion2(T *testing.T) {

	logger.InitLogger("debug")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("连接redis失败, error=" + err.Error())
	}
	rdb.FlushAll(ctx) // 清空 Redis

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		Producer(rdb)
	}()

	go func() {
		defer wg.Done()
		Scheduler(rdb)
	}()

	// 启动两个 Worker，每个有独立的 Stream
	go func() {
		defer wg.Done()
		Worker(rdb, "worker1", workerStreams[0])
	}()

	go func() {
		defer wg.Done()
		Worker(rdb, "worker2", workerStreams[1])
	}()

	wg.Wait()
}

func Worker(rdb *redis.Client, workName string, streamName string) {
	rdb.XGroupCreateMkStream(ctx, streamName, "worker_group", "0")
	for {
		results, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "worker_group",
			Consumer: workName,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    0, //阻塞等待
		}).Result()
		if err != nil {
			logger.Errorf("%s 读取 Stream 失败:%v", workName, err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, result := range results {
			for _, msg := range result.Messages {
				taskId, ok := msg.Values["task_id"].(string)
				if !ok {
					logger.Errorf("%s 流数据获取任务id失败 %v", workName, msg.Values["task_id"])
					continue
				}

				taskJson, err := rdb.HGet(ctx, TaskMetadata, taskId).Result()
				if err != nil {
					logger.Errorf("%s 获取元数据失败 taskId = %v err = %v", workName, taskId, err)
					continue
				}

				task := &Task{}
				if err = json.Unmarshal([]byte(taskJson), task); err != nil {
					logger.Errorf("%s 元数据转josn失败,val = %s,err = %v", workName, taskJson, err)
					continue
				}

				//加锁防止重复处理
				lockKey := fmt.Sprintf("lock:%s", taskId)
				locked, _ := rdb.SetNX(ctx, lockKey, workName, 10*time.Second).Result()
				if !locked {
					continue
				}
				defer rdb.Del(ctx, lockKey)

				//更新状态
				task.Status = "processing"
				task.StartedAt = time.Now().Unix()
				taskJSON, _ := json.Marshal(task)
				if err = rdb.HSet(ctx, TaskMetadata, taskId, taskJSON).Err(); err != nil {
					logger.Errorf("%s将状态更改为处理中失败: %v", workName, err)
				}

				// 模拟执行
				time.Sleep(3 * time.Second)
				logger.Infof("%s 处理任务: %d, 优先级: %d, 数据: %s", workName, task.ID, task.Priority, task.Data)

				if rand.Float32() < 0.5 && task.Priority < 3 { //最多重试3次
					task.RetryCount++
					task.Status = "pending"
					task.ExecuteAt = time.Now().Unix() + 2
					taskJSON, _ = json.Marshal(task)
					rdb.HSet(ctx, TaskMetadata, taskId, taskJSON)
					rdb.ZAdd(ctx, TasksPending, &redis.Z{Score: float64(task.ExecuteAt), Member: taskId})
					logger.Infof("%s 任务 %d 失败，重试 %d", workName, task.ID, task.RetryCount)
				} else if task.RetryCount >= 3 {
					task.Status = "failed"
					taskJSON, _ = json.Marshal(task)
					rdb.HSet(ctx, TaskMetadata, taskId, taskJSON)
					rdb.ZAdd(ctx, TaskQueueFailed, &redis.Z{Score: float64(task.Priority), Member: taskId})
					logger.Infof("%s 任务 %d 重试超限，标记失败", workName, task.ID)
				} else {
					task.Status = "success"
					taskJSON, _ = json.Marshal(task)
					rdb.HSet(ctx, TaskMetadata, taskId, taskJSON)
					logger.Infof("%s 任务 %d 执行成功", workName, task.ID)
				}
				// 确认消息
				rdb.XAck(ctx, streamName, "worker_group", msg.ID)
			}
		}
	}
}

func Scheduler(rdb *redis.Client) {
	workerIndex := 0 // 轮询分配索引
	for {
		//获取最早到期任务
		earliest, err := rdb.ZRangeWithScores(ctx, TasksPending, 0, 0).Result()
		if err != nil {
			logger.Errorf("获取最早到期任务失败: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		now := time.Now().Unix()
		if len(earliest) == 0 {
			logger.Debugf("待调度队列为空，等待 1 秒")
			time.Sleep(1 * time.Second)
			continue
		}
		earliestTime := int64(earliest[0].Score)
		if earliestTime > now {
			waitDuration := time.Duration(earliestTime-now) * time.Second
			logger.Debugf("最早任务 %v 未到期，等待 %v", earliest[0].Member, waitDuration)
			time.Sleep(waitDuration)
		}

		//获取所有任务
		tasks, err := rdb.ZRangeByScore(ctx, TasksPending, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    fmt.Sprintf("%v", now),
			Offset: 0,
			Count:  10,
		}).Result()
		if err != nil {
			logger.Errorf("获取所有任务失败: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		//负载均衡：轮询分配到 Worker Stream
		for _, taskId := range tasks {
			taskJSON, err := rdb.HGet(ctx, TaskMetadata, taskId).Result()
			if err != nil {
				logger.Errorf("获取元数据失败id=%s,err=%v", taskId, err)
				continue
			}

			task := &Task{}
			_ = json.Unmarshal([]byte(taskJSON), task)

			targetStream := workerStreams[workerIndex]
			err = rdb.XAdd(ctx, &redis.XAddArgs{
				Stream: targetStream,
				Values: map[string]interface{}{
					"task_id":   taskId,
					"priority":  task.Priority,
					"task_json": taskJSON,
				},
			}).Err()
			if err != nil {
				logger.Errorf("调度任务 %s 到 %s 失败 %s", taskId, targetStream, err)
			}

			rdb.ZRem(ctx, TasksPending, taskId)
			logger.Infof("调度任务 %s 到 %s", taskId, targetStream)

			workerIndex = (workerIndex + 1) % len(workerStreams) // 轮询下一个 Worker
		}
		time.Sleep(1 * time.Second)
	}
}

// Producer
func Producer(rdb *redis.Client) {
	var ids = []uint64{345144687782801408, 345144821595779072, 345145023033520128, 345145116700717056, 345146225150406656}
	for i, id := range ids {
		task := Task{
			ID:         id,
			Priority:   rand.Intn(10),
			Data:       fmt.Sprintf("task_%d,sort_%d", id, i+1),
			ExecuteAt:  time.Now().Unix() + int64(i+1),
			Status:     "pending",
			RetryCount: 0,
			StartedAt:  0,
		}
		taskJSON, _ := json.Marshal(task)
		rdb.ZAdd(ctx, TasksPending, &redis.Z{Score: float64(task.ExecuteAt), Member: fmt.Sprintf("%d", task.ID)})
		rdb.HSet(ctx, TaskMetadata, fmt.Sprintf("%d", task.ID), taskJSON)
		logger.Infof("produce 生产了数据 %+v", task)
		time.Sleep(1 * time.Second)
	}
}
