package pkg

import "github.com/redis/go-redis/v9"

func CreateCacheStore() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})
}
