package apps

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/sentry/producer"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/minoic/glgf"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"sync"
	"time"
)

func NewApplication(ctx context.Context, deadline time.Time, entity *models.Application) (*Application, error) {
	handler := parseImplements(ctx, entity)
	if handler == nil {
		return nil, errors.New("cant find this application type")
	}
	ctx2, cancel := context.WithDeadline(ctx, deadline)
	app := &Application{
		ctx:         ctx2,
		Cancel:      cancel,
		deadline:    deadline,
		Handler:     handler,
		measurement: common.Measurements[entity.SchemeID],
	}
	err := app.UpdateOptions(entity)
	if err != nil {
		return nil, err
	}
	return app, nil
}

type Application struct {
	ctx         context.Context
	Cancel      func()
	Options     models.Application
	Handler     Implement
	measurement string
	deadline    time.Time
	failedCount uint
	once        sync.Once
	startTime   time.Time
}

func (this *Application) Run() {
	this.once.Do(func() {
		go func() {
			glgf.Info("Application started:", this.Options.Name, this.Options.ID, "till", this.deadline, "rest time:", fmt.Sprintf("%.2f", this.deadline.Sub(time.Now()).Seconds()), "(s)")
			defer glgf.Warn("Application stopped:", this.Options.Name, this.Options.ID)
			t1 := time.NewTicker(5 * time.Second)
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

const THRESHOLD = 3

func (this *Application) refresh() {
	sentryTime := time.Now()
	status, err := this.Handler.Fetch(this.ctx)
	if status == nil {
		panic(err)
	}
	fields := map[string]any{}
	_ = mapstructure.Decode(status, &fields)
	point := write.NewPoint(
		this.measurement,
		map[string]string{
			"id": this.Options.ID.Hex(),
		},
		fields,
		sentryTime,
	)
	if err != nil {
		// 状态不正常时
		point.AddField("c_err", err.Error())
		point.AddField("c_live", false)
		point.AddField("c_start_time", time.Now())
		this.failedCount = min(this.failedCount+1, THRESHOLD+1)
		if this.failedCount < THRESHOLD && (this.failedCount&1 == 0) {
			defer this.reclaim()
		} else if this.failedCount == THRESHOLD {
			this.alert(err.Error())
		}
	} else {
		// 状态正常时
		// 防抖
		if this.failedCount != 0 {
			if this.failedCount >= THRESHOLD {
				this.failedCount = THRESHOLD
			}
			this.failedCount -= THRESHOLD / 3
			if this.failedCount < 0 {
				this.failedCount = 0
			}
			if this.failedCount == 0 {
				// 确认恢复
				this.startTime = time.Now()
				this.onlineLog()
			}
		}
		point.AddField("c_err", "")
		point.AddField("c_live", true)
		point.AddField("c_start_time", this.startTime)
	}
	point.AddField("c_failed_count", this.failedCount)
	point.AddField("c_sentry", common.Hostname())
	storage.WriteInfluxDB.WritePoint(point)
}

func (this *Application) reclaim() {
	storage.WriteInfluxDB.Flush()
	glgf.Debug("reclaiming", this.Options.Name)
	err := producer.EntityClaim(this.ctx, this.Options.ID.Hex(), this.deadline, &this.Options)
	if err != nil {
		glgf.Warn("cant reclaim entity,", err)
		return
	}
	this.Cancel()
}

func (this *Application) onlineLog() {
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

func (this *Application) alert(msg string) {
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

func (this *Application) UpdateOptions(option *models.Application) error {
	this.Options = *option
	err := this.Handler.Init(func(options any) error {
		return mapstructure.Decode(this.Options.Options, options)
	})
	if err != nil {
		return err
	}
	query, err := storage.QueryInfluxDB.Query(
		this.ctx,
		fmt.Sprintf(
			`
from(bucket: "backend")
  |> range(start: -1h)
  |> last()
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_failed_count" or r["_field"] == "c_start_time")
  |> filter(fn: (r) => r["id"] == "%s")`,
			this.measurement,
			this.Options.ID.Hex(),
		),
	)
	if err != nil {
		return err
	}
	for query.Next() {
		if query.Record().Field() == "c_start_time" {
			this.startTime = cast.ToTime(query.Record().Value())
			glgf.Debug("entity start time:", this.startTime)
		} else if query.Record().Field() == "c_failed_count" {
			this.failedCount = cast.ToUint(query.Record().Value())

			glgf.Debug("entity failed count:", this.failedCount)
		}
	}
	if this.startTime.IsZero() {
		this.startTime = time.Now()
	}
	return nil
}
