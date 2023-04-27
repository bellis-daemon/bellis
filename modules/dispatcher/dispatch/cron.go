package dispatch

import (
	"github.com/robfig/cron/v3"
)

func RunTasks() {
	c := cron.New()
	c.AddFunc("@every 10s", syncEntityID)
	c.AddFunc("@every 100ms", checkEntities)
	c.Run()
}
