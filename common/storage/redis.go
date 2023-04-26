package storage

import (
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client

func Redis() *redis.Client {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:        "redis:6379",
			DialTimeout: 3 * time.Second,
		})
	}
	return rdb
}
