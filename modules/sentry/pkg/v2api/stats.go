package v2api

import (
	"context"
	"github.com/v2fly/v2ray-core/v5/app/stats/command"
)

func NodeSysStatus(host string) (*command.SysStatsResponse, error) {
	cmdConn, err := getGrpcConn(host)
	if err != nil {
		return nil, err
	}
	statsClient := command.NewStatsServiceClient(cmdConn)
	return statsClient.GetSysStats(context.Background(), &command.SysStatsRequest{})
}
