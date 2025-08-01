package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestList(T *testing.T) {
	var wg sync.WaitGroup

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("连接redis失败, error=" + err.Error())
	}

	//go ProductionList(rdb)
	ProductionList(rdb, &wg)
	//go Customer(rdb)
	//time.Sleep(1 * time.Second)
	//wg.Wait()
}

func ProductionList(rdb *redis.Client, wg *sync.WaitGroup) {
	wg.Add(1)
	fmt.Println("生产者开始生产数据")
	var err error
	for i := 1; i <= 10; i++ {
		err = rdb.LPush(context.Background(), "task_queue", fmt.Sprintf("%d", i)).Err()
		if err != nil {
			fmt.Println("Lpush Error: ", err.Error())
			return
		}
		fmt.Println("生产者生产了数据: ", i)
		time.Sleep(1 * time.Second)
	}
}
func Customer(rdb *redis.Client, wg *sync.WaitGroup) {
	fmt.Println("消费者开始消费数据")
	defer wg.Done()
	count := 0
	for {
		result, err := rdb.BRPop(context.Background(), 0, "task_queue").Result()
		if err != nil {
			fmt.Println("Pop error : ", err.Error())
			return
		}
		time.Sleep(500 * time.Millisecond)
		fmt.Println("get Value = ", result[1])
		count++
		if count >= 10 {
			fmt.Println("数据读取完了，退出")
			break
		}
	}
}
