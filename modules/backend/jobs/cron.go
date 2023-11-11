package jobs

import (
	"context"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/go-co-op/gocron"
	redislock "github.com/go-co-op/gocron-redis-lock"
	"github.com/minoic/glgf"
	"time"
)

func StartAsync() {
	s := gocron.NewScheduler(time.Local)
	locker, err := redislock.NewRedisLocker(storage.Redis())
	if err != nil {
		panic(err)
	}
	s.WithDistributedLocker(locker)
	s.Every("3h").StartImmediately().Do(func() {
		glgf.Info("deleting data before one week in influxdb")
		err := storage.DeleteInfluxDB.DeleteWithName(context.Background(), "bellis", "backend", time.UnixMilli(0), time.Now().AddDate(0, 0, -7), "")
		if err != nil {
			glgf.Error(err)
		}
	})
	s.StartAsync()
}
