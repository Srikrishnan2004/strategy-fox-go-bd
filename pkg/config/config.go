package config

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"log"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to redis successfully.")
}
