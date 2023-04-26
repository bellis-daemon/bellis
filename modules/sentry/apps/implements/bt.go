package implements

import (
	"context"
	btgosdk "github.com/minoic/bt-go-sdk"
)

type BT struct {
	Options btOptions
	Client  btgosdk.Client
}

func (this *BT) Init(setOptions func(options any) error) error {
	return setOptions(&this.Options)
}

func (this *BT) Fetch(ctx context.Context) (any, error) {
	this.Client.BTAddress = this.Options.Address
	this.Client.BTKey = this.Options.Token
	ret, err := this.Client.GetNetWork()
	if err != nil {
		return &btStatus{}, err
	} else {
		ret.CPU = append(ret.CPU, 0, 0)
		return btStatus{
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
