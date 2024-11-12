package pkg

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func CreateCacheStore() *redis.Client {
	dsn := fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	return redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
		Protocol: 2,
	})
}
