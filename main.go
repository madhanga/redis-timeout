package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Redis addresses: master and replica
	addresses := []string{"localhost:6379", "localhost:6380"}
	currentAddrIndex := 0
	failureCount := 0
	const maxFailures = 3 // Switch after 3 consecutive failures

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addresses[currentAddrIndex],
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
	fmt.Printf("Connected to Redis at %s: %s\n", addresses[currentAddrIndex], pong)

	// Continuously read key in a loop
	keyName := "mykey"
	for {
		// Create a new context for each request with timeout
		reqCtx, reqCancel := context.WithTimeout(context.Background(), 2*time.Second)
		
		val, err := client.Get(reqCtx, keyName).Result()
		timestamp := time.Now().Format("15:04:05")
		
		if err == redis.Nil {
			fmt.Printf("[%s] [%s] Key '%s' does not exist\n", timestamp, addresses[currentAddrIndex], keyName)
			failureCount = 0 // Reset on success
		} else if err != nil {
			fmt.Printf("[%s] [%s] Error reading key: %v\n", timestamp, addresses[currentAddrIndex], err)
			failureCount++
			
			// Switch to the other Redis instance after consecutive failures
			if failureCount >= maxFailures {
				currentAddrIndex = (currentAddrIndex + 1) % len(addresses)
				fmt.Printf("[%s] Switching to %s\n", timestamp, addresses[currentAddrIndex])
				
				// Close old client and create new one
				client.Close()
				client = redis.NewClient(&redis.Options{
					Addr:     addresses[currentAddrIndex],
					Password: "",
					DB:       0,
				})
				failureCount = 0
			}
		} else {
			fmt.Printf("[%s] [%s] Key '%s' = '%s'\n", timestamp, addresses[currentAddrIndex], keyName, val)
			failureCount = 0 // Reset on success
		}
		
		reqCancel()
		time.Sleep(1 * time.Second)
	}
}