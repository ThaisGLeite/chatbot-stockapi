package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func GetRedisClient() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(pong, err)
	}
	return rdb
}
