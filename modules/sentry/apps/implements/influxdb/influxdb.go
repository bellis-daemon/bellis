package influxdb

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
)

type InfluxDB struct {
	options influxDBOptions
}

func (this *InfluxDB) Fetch(ctx context.Context) (status.Status, error) {
	//TODO implement me
	panic("implement me")
}

func (this *InfluxDB) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
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
