package implements

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	btgosdk "github.com/minoic/bt-go-sdk"
)

type BT struct {
	options btOptions
	client  btgosdk.Client
}

func (this *BT) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
}

func (this *BT) Fetch(ctx context.Context) (status.Status, error) {
	this.client.BTAddress = this.options.Address
	this.client.BTKey = this.options.Token
	ret, err := this.client.GetNetWork()
	if err != nil {
		return &btStatus{}, err
	} else {
		ret.CPU = append(ret.CPU, 0, 0)
		return &btStatus{
			MemFree:   ret.Mem.MemFree,
			MemTotal:  ret.Mem.MemTotal,
			Up:        ret.Up,
			Down:      ret.Down,
			UpTotal:   ret.UpTotal,
			DownTotal: ret.DownTotal,
		}, nil
	}
}

type btOptions struct {
	Address string
	Token   string
}

type btStatus struct {
	MemFree   int
	MemTotal  int
	CPUUsage  int
	CPUCores  float64
	Up        float64
	Down      float64
	UpTotal   int64
	DownTotal int64
}

func (this *btStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}
