package v2ray

import (
	"fmt"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/bellis-daemon/bellis/modules/sentry/pkg/v2api"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

type V2ray struct {
	implements.Template
	options v2rayOptions
}

func (this *V2ray) Fetch(ctx context.Context) (status.Status, error) {
	host := fmt.Sprintf("%s:%d", this.options.Address, this.options.Port)
	var err error
	s := &v2rayStatus{}
	s.TagTraffic, err = v2api.NodeTagTraffic(host, this.options.Tag)
	if err != nil {
		return &v2rayStatus{}, err
	}
	stats, err := v2api.NodeSysStatus(host)
	if err != nil {
		return &v2rayStatus{}, err
	}
	s.NumGoroutine = stats.NumGoroutine
	s.NumGC = stats.NumGC
	s.Alloc = stats.Alloc
	s.TotalAlloc = stats.TotalAlloc
	s.Sys = stats.Sys
	s.Mallocs = stats.Mallocs
	s.Frees = stats.Frees
	s.LiveObjects = stats.LiveObjects
	s.PauseTotalNs = stats.PauseTotalNs
	s.Uptime = stats.Uptime
	return s, nil
}

type v2rayStatus struct {
	NumGoroutine uint32
	NumGC        uint32
	Alloc        uint64
	TotalAlloc   uint64
	Sys          uint64
	Mallocs      uint64
	Frees        uint64
	LiveObjects  uint64
	PauseTotalNs uint64
	Uptime       uint32
	TagTraffic   int64
}

func (this *v2rayStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type v2rayOptions struct {
	Address string
	Port    int
	Tag     string
}

func init() {
	implements.Register("v2ray", func(options bson.M) implements.Implement {
		return &V2ray{options: option.ToOption[v2rayOptions](options)}
	})
}
