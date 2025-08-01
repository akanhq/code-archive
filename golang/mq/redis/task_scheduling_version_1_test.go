package main

//import (
//	"code_test/logger"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/go-redis/redis/v8"
//	"math/rand"
//	"sync"
//	"testing"
//	"time"
//)
//
//var ctx = context.Background()
//
//const (
//	TasksPending    = "tasks_pending"     //调度队列(Sorted Set)
//	TaskMetadata    = "task_metadata"     //任务元数据(Hash)
//	TaskQueue       = "task_queue"        //执行队列(Sorted Set)
//	TaskQueueFailed = "task_queue_failed" //失败队列(Sorted Set)-失败次数超过指定次数
//)
//
//// Task 定义任务结构体
//type Task struct {
//	ID         uint64 `json:"id"`
//	Priority   int    `json:"priority"`    // 优先级，越大越优先
//	Data       string `json:"data"`        // 任务内容
//	ExecuteAt  int64  `json:"execute_at"`  // 执行时间戳
//	Status     string `json:"status"`      // 状态: pending, processing, success, failed
//	RetryCount int    `json:"retry_count"` // 重试次数
//	StartedAt  int64  `json:"started_at"`  // 开始处理时间戳（用于超时检测）
//}
//
//func TestTaskSchedulingVersion1(T *testing.T) {
//
//	logger.InitLogger("debug")
//
//	rdb := redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//	})
//	_, err := rdb.Ping(ctx).Result()
//	if err != nil {
//		panic("连接redis失败, error=" + err.Error())
//	}
//
//	var wg sync.WaitGroup
//	wg.Add(3)
//
//	go func() {
//		defer wg.Done()
//		Producer(rdb)
//	}()
//
//	go func() {
//		defer wg.Done()
//		Scheduler(rdb)
//	}()
//
//	go func() {
//		defer wg.Done()
//		Worker(rdb, "work111111---->")
//	}()
//
//	go func() {
//		defer wg.Done()
//		Worker(rdb, "work222222---->")
//	}()
//
//	wg.Wait()
//}
//
//// 提交任务
//func Producer(rdb *redis.Client) {
//
//	var ids = []uint64{345144687782801408,
//		345144821595779072,
//		345145023033520128,
//		345145116700717056,
//		345146225150406656,
//		345152428637171712,
//		345248430343860224,
//		345254056639868928,
//		345254849531097088}
//
//	for i, id := range ids {
//		task := Task{
//			ID:         id,
//			Priority:   rand.Intn(10),                           //优先级
//			Data:       fmt.Sprintf("task_%d,sort_%d", id, i+1), //处理内容
//			ExecuteAt:  time.Now().Unix(),                       //执行时间
//			Status:     "pending",                               //处理中
//			RetryCount: 0,                                       //重试次数
//			StartedAt:  0,                                       //开始处理时间戳
//		}
//
//		//添加到待处理集合，有序集合
//		err := rdb.ZAdd(ctx, TasksPending, &redis.Z{
//			Score:  float64(task.ExecuteAt),
//			Member: task.ID,
//		}).Err()
//		if err != nil {
//			logger.Errorf("添加到待处理集合失败 %d: %v", id, err)
//			continue
//		}
//
//		//将元数据添加到hash
//		taskJson, err := json.Marshal(&task)
//		if err != nil {
//			logger.Errorf("task to json error: %v", err)
//		}
//
//		err = rdb.HSet(ctx, TaskMetadata, fmt.Sprintf("%d", task.ID), taskJson).Err()
//		if err != nil {
//			logger.Errorf("将元数据添加到hash error: %v", err)
//		}
//		logger.Infof("produce 生产了数据%+v", task)
//		time.Sleep(1 * time.Second)
//	}
//}
//
//// Scheduler 调度任务
//func Scheduler(rdb *redis.Client) {
//	for {
//		//获取时间小于现在的任务
//		tasks, err := rdb.ZRangeByScore(ctx, TasksPending, &redis.ZRangeBy{
//			Min:    "-inf",
//			Max:    fmt.Sprintf("%v", time.Now().Unix()),
//			Offset: 0,
//			Count:  10,
//		}).Result()
//		if err != nil {
//			logger.Infof("获取时间小于现在的任务失败%v", err)
//			//一秒后重试
//			time.Sleep(1 * time.Second)
//		}
//
//		for _, taskId := range tasks {
//			// 获取任务优先级
//			taskJSON, err := rdb.HGet(ctx, TaskMetadata, taskId).Result()
//			if err != nil {
//				logger.Errorf("获取任务优先级失败hash error: %s", err)
//			}
//
//			var task Task
//			json.Unmarshal([]byte(taskJSON), &task)
//
//			//添加数据到执行队列
//			rdb.ZAdd(ctx, TaskQueue, &redis.Z{
//				Score:  float64(task.Priority),
//				Member: taskId,
//			})
//			//	在pending中删除待处理数据
//			err = rdb.ZRem(ctx, TasksPending, taskId).Err()
//			if err != nil {
//				logger.Errorf("从pending中删除数据失败：%v", err)
//			}
//			time.Sleep(500 * time.Millisecond)
//		}
//		time.Sleep(1 * time.Second)
//	}
//}
//
//// Worker 执行任务
//func Worker(rdb *redis.Client, workName string) {
//	for {
//		taskIds, err := rdb.ZRevRange(ctx, TaskQueue, 0, 0).Result() //取最高优先级
//		if err != nil || len(taskIds) == 0 {
//			logger.Warnf("队列为空，短暂等待")
//			time.Sleep(500 * time.Millisecond)
//			continue
//		}
//		taskId := taskIds[0]
//
//		//拿到元数据
//		taskJson, err := rdb.HGet(ctx, TaskMetadata, taskId).Result()
//		if err != nil {
//			logger.Errorf("%s获取任务元数据失败: %v", workName, err)
//			time.Sleep(1 * time.Second)
//			continue
//		}
//		var task Task
//		json.Unmarshal([]byte(taskJson), &task)
//
//		logger.Infof("%s获取任务元数据: %+v", workName, task)
//
//		//元数据的状态在等待处理中才能处理pending
//		if task.Status == "processing" {
//			logger.Infof("%s 正在处理", taskId)
//			time.Sleep(3 * time.Second)
//			continue
//		}
//
//		//使用分布式锁，确保同一时间只有一个worker处理任务
//		lockKey := fmt.Sprintf("lock:%s", taskId)
//		locked, err := rdb.SetNX(ctx, lockKey, workName, 10*time.Second).Result()
//		if err != nil || !locked {
//
//			//这里有一个问题，就是多个协程同时拿到元数据，但是有一个协程获取到锁，拿元数据修改后，另外一个协程就获取锁失败，打印的状态是修改前的元数据，并不是修改后的元数据
//			logger.Errorf("%s获取锁失败,lockKey=%s,status=%s,err=%v", workName, lockKey, task.Status, err)
//			time.Sleep(1 * time.Second)
//			continue
//		}
//		defer rdb.Del(ctx, lockKey)
//
//		logger.Infof("%s获取锁成功,lockKey=%s", workName, lockKey)
//
//		//将状态更改为处理中
//		task.Status = "processing"
//		taskJSON, err := json.Marshal(task)
//		if err != nil {
//			logger.Errorf("%stask to json error: %v", workName, err)
//			time.Sleep(1 * time.Second)
//			continue
//		}
//		err = rdb.HSet(ctx, TaskMetadata, taskId, taskJSON).Err()
//		if err != nil {
//			logger.Errorf("%s将状态更改为处理中失败: %v", workName, err)
//			time.Sleep(1 * time.Second)
//			continue
//		}
//
//		//再次获取修改后的元数据
//		taskJson, err = rdb.HGet(ctx, TaskMetadata, taskId).Result()
//		if err != nil {
//			logger.Errorf("%s获取任务元数据失败: %v", workName, err)
//			time.Sleep(1 * time.Second)
//			continue
//		}
//		var tempTask Task
//		json.Unmarshal([]byte(taskJson), &tempTask)
//		logger.Infof("%s修改后的状态: %v,taskId=%s", workName, tempTask.Status, taskId)
//
//		// 模拟任务执行
//		time.Sleep(3 * time.Second)
//		fmt.Printf("%s任务正在运行: %d, 优先级: %d, 数据: %s\n", workName, task.ID, task.Priority, task.Data)
//
//		//模拟执行结果，50%的失败率，最多重试3次
//		if rand.Float32() < 0.5 && task.Priority < 3 { //最多重试3次
//
//			//处理失败，重新添加到待处理列表
//			task.RetryCount++
//			task.Status = "pending"
//			task.ExecuteAt = time.Now().Unix() + 2 //两秒后重试
//			taskJSON, err = json.Marshal(task)
//			if err != nil {
//				logger.Errorf("task to json error: %v", err)
//			}
//			err = rdb.HSet(ctx, TaskMetadata, taskId, taskJSON).Err()
//			if err != nil {
//				logger.Errorf("%s将状态更改为处理中失败: %v", workName, err)
//			}
//
//			err = rdb.ZAdd(ctx, TasksPending, &redis.Z{
//				Score:  float64(task.ExecuteAt),
//				Member: task.ID,
//			}).Err()
//			if err != nil {
//				logger.Errorf("%s任务添加到重试失败: %v", workName, err)
//			}
//
//			logger.Debugf("%s%d:任务处理失败进行重试,重试次数:%d", workName, task.ID, task.RetryCount)
//		} else if task.RetryCount >= 3 {
//			task.Status = "failed"
//			taskJSON, err = json.Marshal(task)
//			if err != nil {
//				logger.Errorf("%stask to json error: %v", workName, err)
//			}
//			err = rdb.HSet(ctx, TaskMetadata, taskId, taskJSON).Err()
//			if err != nil {
//				logger.Infof("%s任务重试次数超限，已标记为失败: %d", workName, task.ID)
//			}
//
//			//添加到失败队列
//			err = rdb.ZAdd(ctx, TaskQueueFailed, &redis.Z{
//				Score:  float64(task.Priority),
//				Member: taskId,
//			}).Err()
//			if err == nil {
//				//从执行队列中删除
//				err = rdb.ZRem(ctx, TaskQueue, taskId).Err()
//				if err != nil {
//					time.Sleep(500 * time.Millisecond)
//					err = rdb.ZRem(ctx, TaskQueue, taskId).Err()
//					if err != nil {
//						logger.Errorf("从执行队列中删除数据失败%v", err)
//					}
//				}
//			} else {
//				logger.Errorf("添加到失败队列失败")
//			}
//
//		} else {
//			//执行成功，将元数据中的状态改为已成功
//			task.Status = "success"
//			taskJSON, err = json.Marshal(task)
//			if err != nil {
//				logger.Errorf("%stask to json error: %v", workName, err)
//			}
//			err = rdb.HSet(ctx, TaskMetadata, taskId, taskJSON).Err()
//			if err != nil {
//				logger.Errorf("%s%d将状态更改为成功失败: %v", workName, task.ID, err)
//			}
//
//			//从执行队列中删除数据
//			err = rdb.ZRem(ctx, TaskQueue, taskId).Err()
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				err = rdb.ZRem(ctx, TaskQueue, taskId).Err()
//				if err != nil {
//					logger.Errorf("从执行队列中删除数据失败%v", err)
//				}
//			}
//			logger.Infof("%s任务执行成功%d,执行次数%d", workName, task.ID, task.RetryCount)
//		}
//		time.Sleep(2 * time.Second) // 模拟耗时
//	}
//}
