package ping

// You must use pinger.SetPrivileged(true), otherwise you will receive an error

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/go-ping/ping"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Ping struct {
	implements.Template
	options pingOptions
}

type pingOptions struct {
	Address       string
	LossThreshold float64
}

type pingStatus struct {
	PacketLoss float64
	MaxRtt     float64
	MinRtt     float64
	AvgRtt     float64
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
	client.Interval = time.Millisecond
	client.Count = 100
	err = client.Run()
	defer client.Stop()
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
		MaxRtt:     float64(statistics.MaxRtt.Microseconds()) / 1000.0,
		MinRtt:     float64(statistics.MinRtt.Microseconds()) / 1000.0,
		AvgRtt:     float64(statistics.AvgRtt.Microseconds()) / 1000.0,
		IP:         statistics.IPAddr.IP.String(),
	}, nil
}

func init() {
	implements.Register("ping", func(options bson.M) implements.Implement {
		return &Ping{options: option.ToOption[pingOptions](options)}
	})
}
