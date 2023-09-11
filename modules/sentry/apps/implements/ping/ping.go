package ping

// You must use pinger.SetPrivileged(true), otherwise you will receive an error

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/go-ping/ping"
	"github.com/spf13/cast"
	"time"
)

type Ping struct {
	options pingOptions
}

func (this *Ping) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
}

type pingOptions struct {
	Address       string
	LossThreshold float64
}

type pingStatus struct {
	PacketLoss float64
	MaxRtt     int64
	MinRtt     int64
	AvgRtt     int64
	IP         string
}

func (this *pingStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

func (this *Ping) Fetch(ctx context.Context) (status.Status, error) {
	client, err := ping.NewPinger(this.options.Address)
	defer client.Stop()
	if err != nil {
		return &pingStatus{}, err
	}
	client.SetPrivileged(true)
	client.Timeout = 3 * time.Second
	client.Interval = 100 * time.Millisecond
	client.Count = 3
	err = client.Run()
	if err != nil {
		return &pingStatus{}, err
	}
	statistics := client.Statistics()
	if statistics.PacketsRecv == 0 {
		return &pingStatus{}, errors.New("Address " + statistics.IPAddr.IP.String() + " unreachable from our server")
	} else if statistics.PacketLoss > this.options.LossThreshold {
		return &pingStatus{}, errors.New("Address " + statistics.IPAddr.IP.String() + " Excessive packet loss from our server:" + cast.ToString(statistics.PacketLoss))
	}
	return &pingStatus{
		PacketLoss: statistics.PacketLoss,
		MaxRtt:     statistics.MaxRtt.Milliseconds(),
		MinRtt:     statistics.MinRtt.Milliseconds(),
		AvgRtt:     statistics.AvgRtt.Milliseconds(),
		IP:         statistics.IPAddr.IP.String(),
	}, nil
}

func init() {
	implements.Add("ping", func() implements.Implement {
		return &Ping{}
	})
}
