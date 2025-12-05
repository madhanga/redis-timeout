package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping to test connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)

	// Continuously read key in a loop
	keyName := "mykey"
	for {
		// Create a new context for each request with timeout
		reqCtx, reqCancel := context.WithTimeout(context.Background(), 2*time.Second)

		val, err := client.Get(reqCtx, keyName).Result()
		timestamp := time.Now().Format("15:04:05")

		if err == redis.Nil {
			fmt.Printf("[%s] Key '%s' does not exist\n", timestamp, keyName)
		} else if err != nil {
			fmt.Printf("[%s] Error reading key: %v\n", timestamp, err)
		} else {
			fmt.Printf("[%s] Key '%s' = '%s'\n", timestamp, keyName, val)
		}

		reqCancel()
		time.Sleep(1 * time.Second)
	}
}
