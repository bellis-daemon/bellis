package implements

import (
	"fmt"
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
	status := &v2rayStatus{}
	status.TagTraffic, err = v2api.NodeTagTraffic(this.host, this.options.Tag)
	if err != nil {
		return &v2rayStatus{}, err
	}
	stats, err := v2api.NodeSysStatus(this.host)
	if err != nil {
		return &v2rayStatus{}, err
	}
	status.NumGoroutine = stats.NumGoroutine
	status.NumGC = stats.NumGC
	status.Alloc = stats.Alloc
	status.TotalAlloc = stats.TotalAlloc
	status.Sys = stats.Sys
	status.Mallocs = stats.Mallocs
	status.Frees = stats.Frees
	status.LiveObjects = stats.LiveObjects
	status.PauseTotalNs = stats.PauseTotalNs
	status.Uptime = stats.Uptime
	return status, nil
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
