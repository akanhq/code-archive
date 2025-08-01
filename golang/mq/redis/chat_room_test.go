package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var chatRoom = "chat_room"

func TestChatRoom(T *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("连接redis失败, error=" + err.Error())
	}

	go Published(rdb)
	Subscribers(rdb)
}

func Published(rdb *redis.Client) {
	for i := 1; i <= 10; i++ {
		err := rdb.Publish(context.Background(), chatRoom, fmt.Sprintf("我现在往%s发送了一条消息内容是:%d", chatRoom, i)).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("发送消息失败：err:%s", err.Error()))
			return
		}
		time.Sleep(1 * time.Second)
	}
}
func Subscribers(rdb *redis.Client) {
	subscribe := rdb.Subscribe(context.Background(), chatRoom)
	defer subscribe.Close()

	for message := range subscribe.Channel() {
		fmt.Println(fmt.Sprintf("我在%s,接收了一条数据------------------>%s", message.Channel, message.Payload))
	}
}
