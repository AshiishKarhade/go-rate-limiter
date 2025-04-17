package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

func InitRedisConnection() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
	}
	fmt.Println("Connected to Redis server")
	return client
}
