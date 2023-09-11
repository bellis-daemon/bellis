package v2ray

import (
	"fmt"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/bellis-daemon/bellis/modules/sentry/pkg/v2api"
	"golang.org/x/net/context"
)

type V2ray struct {
	options v2rayOptions
	host    string
}

func (this *V2ray) Fetch(ctx context.Context) (status.Status, error) {
	var err error
	s := &v2rayStatus{}
	s.TagTraffic, err = v2api.NodeTagTraffic(this.host, this.options.Tag)
	if err != nil {
		return &v2rayStatus{}, err
	}
	stats, err := v2api.NodeSysStatus(this.host)
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

func (this *V2ray) Init(setOptions func(options any) error) error {
	err := setOptions(&this.options)
	if err != nil {
		return err
	}
	this.host = fmt.Sprintf("%s:%d", this.options.Address, this.options.Port)
	return nil
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
	implements.Add("v2ray", func() implements.Implement {
		return &V2ray{}
	})
}
