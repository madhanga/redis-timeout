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
		Addr:         "amx-arom-redis0.gpx.uat.angelone.in:6379", //order redis uat9
		Password:     "aromamx$123",                              // no password set
		DB:           0,                                          // use default DB
		PoolTimeout:  4000 * time.Millisecond,
		ReadTimeout:  5000 * time.Millisecond,
		WriteTimeout: 5000 * time.Millisecond,
	})
	defer client.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Ping to test connection
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)

	// Continuously read key in a loop
	keyName := "myKey"
	for {
		// Create a new context for each request with timeout
		reqCtx, reqCancel := context.WithTimeout(context.Background(), 3*time.Second)

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
		time.Sleep(1 * time.Millisecond)
	}
}
