package apps

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"

	"github.com/avast/retry-go/v4"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/bellis-daemon/bellis/modules/sentry/producer"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/minoic/glgf"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

var pool = sync.Pool{
	New: func() any {
		return new(Entity)
	},
}

func NewEntity(ctx context.Context, deadline time.Time, entity *models.Application) (*Entity, error) {
	app := pool.Get().(*Entity)
	ctx2, cancel := context.WithDeadline(ctx, deadline)
	app.ctx = ctx2
	app.cancel = cancel
	app.deadline = deadline
	app.measurement = entity.Scheme
	app.once = sync.Once{}
	err := app.UpdateOptions(entity)
	if err != nil {
		return nil, err
	}
	return app, nil
}

type Entity struct {
	ctx         context.Context
	cancel      func()
	Options     models.Application
	Handler     implements.Implement
	measurement string
	deadline    time.Time
	failedCount int
	once        sync.Once
	threshold   int
}

func (this *Entity) Cancel() {
	if this.cancel != nil {
		this.cancel()
	}
}

func (this *Entity) Run() {
	this.once.Do(func() {
		go func() {
			glgf.Info("Entity started:", this.Options.Name, this.Options.ID, "till", this.deadline, "rest time:", fmt.Sprintf("%.2f", this.deadline.Sub(time.Now()).Seconds()), "(s)")
			defer glgf.Warn("Entity stopped:", this.Options.Name, this.Options.ID)
			go this.refresh()
			multiplier1 := this.Options.Public.Multiplier
			if multiplier1 <= 0 {
				multiplier1 = 1
			}
			multiplier2 := this.Handler.Multiplier()
			if multiplier2 <= 0 {
				multiplier2 = 1
			}
			t1 := time.NewTicker(time.Duration(5*multiplier1*multiplier2) * time.Second)
			for {
				select {
				case <-t1.C:
					go this.refresh()
				case <-this.ctx.Done():
					return
				}
			}
		}()
	})
}

func (this *Entity) saveFetch() (s status.Status, err error) {
	defer func() {
		if r := recover(); r != nil {
			s = nil
			err = errors.New(fmt.Sprint(r))
		}
	}()
	return this.Handler.Fetch(this.ctx)
}

func (this *Entity) refresh() {
	sentryTime := time.Now()
	s, err := this.saveFetch()
	responseTime := time.Since(sentryTime)
	fields := map[string]any{}
	_ = mapstructure.Decode(s, &fields)
	point := write.NewPoint(
		this.measurement,
		map[string]string{
			"id": this.Options.ID.Hex(),
		},
		fields,
		sentryTime,
	)
	cErr := ""
	cLive := true
	if err != nil {
		// err occured, entity offline
		cErr = err.Error()
		cLive = false
		// todo: fix offline judge method
		this.failedCount = min(this.failedCount+1, this.threshold+1)
		if this.failedCount < this.threshold && (this.failedCount&1 == 0) {
			defer this.reclaim()
		} else if this.failedCount == this.threshold {
			this.alert(err.Error())
		}
	} else {
		// no error, entity online
		// debounce
		if this.failedCount != 0 {
			if this.failedCount >= this.threshold {
				this.failedCount = this.threshold
			}
			this.failedCount -= this.threshold / 3
			if this.failedCount < 0 {
				this.failedCount = 0
			}
			if this.failedCount == 0 {
				// confirm that server goes online
				// NOTICE: cant make sure if entity reached failed count threshold before
				this.onlineLog()
			}
		}
		// test the triggers
		for i := range this.Options.Public.TriggerList {
			result := s.PullTrigger(this.Options.Public.TriggerList[i])
			if result != nil {
				this.triggerAlert(result)
			}
		}
	}
	point.AddField("c_err", cErr)
	point.AddField("c_live", cLive)
	point.AddField("c_failed_count", cast.ToUint32(this.failedCount))
	point.AddField("c_sentry", common.Hostname())
	point.AddField("c_response_time", responseTime.Milliseconds())
	storage.WriteInfluxDB.WritePoint(point)
}

func (this *Entity) reclaim() {
	storage.WriteInfluxDB.Flush()
	glgf.Debug("reclaiming", this.Options.Name)
	err := producer.EntityClaim(this.ctx, this.Options.ID.Hex(), this.deadline, &this.Options)
	if err != nil {
		glgf.Warn("cant reclaim entity,", err)
		return
	}
	this.Cancel()
}

func (this *Entity) onlineLog() {
	storage.WriteInfluxDB.Flush()
	time.Sleep(2 * time.Second)
	onlineTime := time.Now()
	glgf.Debug("online logging", this.Options.Name)
	err := retry.Do(func() error {
		return producer.EntityOnline(this.ctx, this.Options.ID.Hex(), onlineTime)
	}, retry.Context(this.ctx), retry.Delay(300*time.Millisecond))
	if err != nil {
		glgf.Error(err)
	}
}

func (this *Entity) alert(msg string) {
	storage.WriteInfluxDB.Flush()
	time.Sleep(2 * time.Second)
	offlineTime := time.Now()
	glgf.Debug("alerting", this.Options.Name)
	err := retry.Do(func() error {
		return producer.EntityOffline(this.ctx, this.Options.ID.Hex(), msg, offlineTime)
	}, retry.Context(this.ctx), retry.Delay(300*time.Millisecond))
	if err != nil {
		glgf.Error(err)
	}
}

func (this *Entity) triggerAlert(info *status.TriggerInfo) {
	// todo: implement function
}

func (this *Entity) UpdateOptions(option *models.Application) (err error) {
	this.Options = *option
	this.Handler, err = implements.Spawn(&this.Options)
	if err != nil {
		return err
	}
	this.threshold = option.Public.Threshold
	if this.threshold == 0 {
		this.threshold = 5
	}
	query, err := storage.QueryInfluxDB.Query(
		this.ctx,
		fmt.Sprintf(
			`
from(bucket: "backend")
  |> range(start: -1m)
  |> last()
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_failed_count")
  |> filter(fn: (r) => r["id"] == "%s")`,
			this.measurement,
			this.Options.ID.Hex(),
		),
	)
	if err != nil {
		return err
	}
	for query.Next() {
		if query.Record().Field() == "c_failed_count" {
			this.failedCount = cast.ToInt(query.Record().Value())
		}
	}
	return nil
}
