package jobs

import (
	"time"

	"github.com/bellis-daemon/bellis/common/storage"
	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
)

func StartAsync() {
	locker, err := redislock.NewRedisLocker(storage.Redis())
	if err != nil {
		panic(err)
	}
	s, err := gocron.NewScheduler(gocron.WithDistributedLocker(locker), gocron.WithLocation(time.Local))
	if err != nil {
		panic(err)
	}
	s.NewJob(gocron.DurationJob(12*time.Hour), gocron.NewTask(setTelegramWebhook), gocron.WithStartAt(gocron.WithStartImmediately()))
	s.NewJob(gocron.DurationRandomJob(10*time.Hour, 14*time.Hour), gocron.NewTask(checkUserLevelExpire))
	s.NewJob(gocron.DurationRandomJob(10*time.Hour, 14*time.Hour), gocron.NewTask(checkUserEntityUsageCount))
	s.NewJob(gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(1), gocron.NewAtTimes(gocron.NewAtTime(6, 0, 0))), gocron.NewTask(resetUserUsages))
	s.Start()
}
