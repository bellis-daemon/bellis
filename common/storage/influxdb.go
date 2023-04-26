package storage

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var WriteInfluxDB api.WriteAPI
var QueryInfluxDB api.QueryAPI

func init() {
	client := influxdb2.NewClientWithOptions(
		"http://influxdb:8086",
		"fWk-Z4OYwupoH0N8XoGehjWeI2smqOfsHmZ_SXdwUunG-4xjpjB8iD_WKvEIaMYbE_6fCbujo3l7USv5lxm5DQ==",
		influxdb2.DefaultOptions().SetBatchSize(50).
			SetFlushInterval(200).
			SetUseGZip(true).
			SetRetryInterval(200).
			SetMaxRetries(3).
			SetMaxRetryInterval(500),
	)
	WriteInfluxDB = client.WriteAPI("bellis", "backend")
	QueryInfluxDB = client.QueryAPI("bellis")
}
