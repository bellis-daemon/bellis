package influxdb

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"go.mongodb.org/mongo-driver/bson"
)

type InfluxDB struct {
	implements.Template
	options influxDBOptions
}

func (this *InfluxDB) Fetch(ctx context.Context) (status.Status, error) {
	//TODO implement me
	panic("implement me")
}

type influxDBOptions struct {
}

type influxDBStatus struct {
}

func (this *influxDBStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

func init() {
	implements.Register("influxdb", func(options bson.M) implements.Implement {
		return &InfluxDB{options: option.ToOption[influxDBOptions](options)}
	})
}
