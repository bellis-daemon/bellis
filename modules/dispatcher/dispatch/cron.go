package dispatch

import (
	"github.com/robfig/cron/v3"
	"time"
)

func RunTasks() {
	t := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-t.C:
				checkEntities()
			}
		}
	}()
	c := cron.New()
	c.AddFunc("@every 10s", syncEntityID)
	c.Run()
}
