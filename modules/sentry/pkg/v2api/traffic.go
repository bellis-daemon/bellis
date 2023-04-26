package v2api

import (
	"context"
	statsService "github.com/v2fly/v2ray-core/v5/app/stats/command"
)

func NodeUserTraffic(host string, email string) (int64, error) {
	// need to set user Email
	return traffic(host, `user>>>`+email, false)
}

func NodeTagTraffic(host string, tag string) (int64, error) {
	// need to set Node inbound tag to "proxy"
	return traffic(host, `inbound>>>`+tag, false)
}

func traffic(host string, request string, reset bool) (int64, error) {
	var (
		up   int64
		down int64
	)
	cmdConn, err := getGrpcConn(host)
	if err != nil {
		return 0, err
	}
	statsClient := statsService.NewStatsServiceClient(cmdConn)
	r := &statsService.GetStatsRequest{
		Name:   request + `>>>traffic>>>uplink`,
		Reset_: reset,
	}
	resp, err := statsClient.GetStats(context.Background(), r)
	if err != nil {
		up = 0
		return 0, err
	} else {
		up = resp.Stat.GetValue()
	}
	r = &statsService.GetStatsRequest{
		Name:   request + `>>>traffic>>>downlink`,
		Reset_: reset,
	}
	resp, err = statsClient.GetStats(context.Background(), r)
	if err != nil {
		down = 0
		return 0, err
	} else {
		down = resp.Stat.GetValue()
	}
	return up + down, nil
}
