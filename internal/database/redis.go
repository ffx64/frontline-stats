package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	ctx := context.Background()

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "6379"
	}

	if password == "" {
		password = ""
	}

	var err error

	for i := 0; i < 10; i++ {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password, // no password set
			DB:       0,        // use default DB
		})

		_, err = client.Ping(ctx).Result()
		if err != nil {
			log.Printf("[redis] error: %v", err)
			continue
		}

		log.Printf("[redis] is running:")
		log.Printf("[redis] client created")
		return client
	}

	log.Printf("[redis] failed to connect")
	return nil
}

func CloseRedisClient(client *redis.Client) error {
	err := client.Close()
	if err != nil {
		return err
	}
	log.Println("[redis] client closed")
	return nil
}
