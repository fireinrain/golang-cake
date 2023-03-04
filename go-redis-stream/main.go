package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPass := os.Getenv("REDIS_PASS")

	ctx := context.Background()
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: redisPass, // no password set
		DB:       0,         // use default DB
	})
	defer rdb.Close()

	// 向stream中添加记录
	id, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "mystream",
		ID:     "*",
		Values: map[string]interface{}{"name": "Alice", "age": "25"},
	}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Added record to mystream with ID", id)

	// 读取stream中的记录
	res, err := rdb.XRange(ctx, "mystream", "-", "+").Result()
	if err != nil {
		panic(err)
	}
	for _, rec := range res {
		fmt.Printf("Record %s: %v\n", rec.ID, rec.Values)
	}

	// 创建消费者组
	consumer := "myconsumer"
	info, err := rdb.XInfoGroups(ctx, "mystream").Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	for _, group := range info {
		if group.Name == consumer {
			fmt.Printf("Consumer group %s exists.\n", consumer)
			//有消费组 先删除
			if _, err := rdb.XGroupDelConsumer(ctx, "mystream", consumer, consumer).Result(); err != nil {
				panic(err)
			}
			break
		}
	}
	fmt.Printf("Consumer group %s does not exist.\n", consumer)

	if _, err := rdb.XGroupCreate(ctx, "mystream", consumer, "0").Result(); err != nil && err != redis.Nil {
		panic(err)
	}

	// 订阅stream
	streams := []string{"mystream", ">"}
	for {
		res, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumer,
			Consumer: consumer,
			Streams:  streams,
			Count:    1,
			Block:    time.Second * 10,
		}).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			panic(err)
		}
		for _, stream := range res {
			for _, msg := range stream.Messages {
				fmt.Printf("Consumer %s received message %s: %v\n", consumer, msg.ID, msg.Values)
				// 处理消息...
				rdb.XAck(ctx, "mystream", consumer, msg.ID)
			}
		}
	}
}
