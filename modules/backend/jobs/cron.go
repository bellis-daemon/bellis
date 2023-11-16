package jobs

import (
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/go-co-op/gocron"
	redislock "github.com/go-co-op/gocron-redis-lock"
	"time"
)

func StartAsync() {
	s := gocron.NewScheduler(time.Local)
	locker, err := redislock.NewRedisLocker(storage.Redis())
	if err != nil {
		panic(err)
	}
	s.WithDistributedLocker(locker)
	s.Every("3h").StartImmediately().Do(clearInfluxdbDataExpired)
	s.Every("12h").StartImmediately().Do(setTelegramWebhook)
	s.StartAsync()
}
