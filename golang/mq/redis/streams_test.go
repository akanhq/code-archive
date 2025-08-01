package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestStreams(T *testing.T) {

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("连接redis失败, error=" + err.Error())
	}

	//创建消费组
	rdb.XGroupCreateMkStream(ctx, "orders", "order_group", "0")

	//生产者
	go func() {
		for i := 1; i <= 10; i++ {
			err = rdb.XAdd(ctx, &redis.XAddArgs{
				Stream: "orders",
				Values: map[string]interface{}{
					"order_id": i,
					"amount":   i * 10,
				},
			}).Err()

			if err != nil {
				fmt.Println(fmt.Sprintf("生产订单数据失败:%s", err))
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	//消费者
	for {
		result, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "order_group",
			Consumer: "consumer1",
			Streams:  []string{"orders", ">"},
			Count:    1,
			Block:    0,
		}).Result()
		if err != nil {
			fmt.Println(fmt.Sprintf("消费数据失败:%s", err))
			return
		}

		for _, message := range result[0].Messages {
			fmt.Println("Processing order:", message.Values)
			rdb.XAck(ctx, "orders", "order_group", message.ID)
		}
	}
}
