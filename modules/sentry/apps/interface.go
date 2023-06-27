package apps

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/producer"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/minoic/glgf"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"sync"
	"time"
)

// Implement 必须在每个子类中实现的实际方法
type Implement interface {
	// Fetch must return non nil status value, or it will panic
	// return error if entity is offline
	Fetch(ctx context.Context) (any, error)
	Init(setOptions func(options any) error) error
}

func NewApplication(ctx context.Context, deadline time.Time, entity models.Application) (*Application, error) {
	var (
		handler Implement
	)
	switch entity.SchemeID {
	case BT:
		handler = &implements.BT{}
	case Ping:
		handler = &implements.Ping{}
	case HTTP:
		handler = &implements.HTTP{}
	case Minecraft:
		handler = &implements.Minecraft{}
	case V2Ray:
		handler = &implements.Minecraft{}
	case DNS:
		handler = &implements.DNS{}
	case VPS:
		handler = &implements.VPS{}
	case Docker:
		handler = &implements.Docker{}
	case Source:
		handler = &implements.Source{}
	default:
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

func (this *Application) refresh() {
	t := time.Now()
	status, err := this.Handler.Fetch(this.ctx)
	if status == nil {
		panic(err)
	}
	m := map[string]any{}
	_ = mapstructure.Decode(status, &m)
	point := write.NewPoint(
		this.measurement,
		map[string]string{
			"id": this.Options.ID.Hex(),
		},
		m,
		t,
	)
	glgf.Debug(status, err)
	if err != nil {
		point.AddField("c_err", err.Error())
		point.AddField("c_live", false)
		this.failedCount++
		if this.failedCount == 1 {
			defer this.reclaim()
		} else if this.failedCount == 2 {
			this.alert(err.Error())
		}
	} else {
		this.failedCount = 0
		point.AddField("c_err", "")
		point.AddField("c_live", true)
	}
	point.AddField("c_failed_count", this.failedCount)
	storage.WriteInfluxDB.WritePoint(point)
}

func (this *Application) reclaim() {
	storage.WriteInfluxDB.Flush()
	glgf.Debug("reclaiming", this.Options.Name)
	err := producer.EntityClaim(this.ctx, this.Options.ID.Hex(), this.deadline)
	if err != nil {
		glgf.Warn("cant reclaim entity,", err)
		return
	}
	this.Cancel()
}

func (this *Application) alert(msg string) {
	retry.Do(func() error {
		glgf.Debug("alerting", this.Options.Name)
		return storage.Redis().XAdd(this.ctx, &redis.XAddArgs{
			Stream: "EntityOfflineAlert",
			MaxLen: 256,
			Approx: true,
			Values: map[string]interface{}{
				"EntityID": this.Options.ID.Hex(),
				"Message":  msg,
			},
		}).Err()
	}, retry.Context(this.ctx), retry.Delay(300*time.Millisecond))
}

func (this *Application) UpdateOptions(option models.Application) error {
	this.Options = option
	err := this.Handler.Init(func(options any) error {
		return mapstructure.Decode(this.Options.Options, options)
	})
	if err != nil {
		return err
	}
	query, err := storage.QueryInfluxDB.Query(
		this.ctx,
		fmt.Sprintf(
			`from(bucket: "backend")
  |> range(start: -1h)
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
	if query.Next() {
		this.failedCount = cast.ToUint(query.Record().Value())
		glgf.Debug("entity failed count:", this.failedCount)
	} else {
		glgf.Warn("cant find failed count in influxdb for entity", this.Options.ID.Hex())
		this.failedCount = 0
	}
	return nil
}
