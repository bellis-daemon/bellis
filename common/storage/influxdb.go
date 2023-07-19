package storage

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var WriteInfluxDB api.WriteAPI
var QueryInfluxDB api.QueryAPI
var DeleteInfluxDB api.DeleteAPI

func init() {
	client := influxdb2.NewClientWithOptions(
		"http://influxdb:8086",
		"nhIDF8c0tM_6dD2ESKy9aqxoPSuzUhAa3onvlUv0h0cXlovRLv6Szfp3TfHzdFXc8emSlw4Mq7T-TkI3fT2uQw==",
		influxdb2.DefaultOptions().SetBatchSize(50).
			SetFlushInterval(200).
			SetUseGZip(true).
			SetRetryInterval(200).
			SetMaxRetries(3).
			SetMaxRetryInterval(500),
	)
	WriteInfluxDB = client.WriteAPI("bellis", "backend")
	QueryInfluxDB = client.QueryAPI("bellis")
	DeleteInfluxDB = client.DeleteAPI()
}
