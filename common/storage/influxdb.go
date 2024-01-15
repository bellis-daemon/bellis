package storage

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var WriteInfluxDB api.WriteAPI
var WriteInfluxDBBlocking api.WriteAPIBlocking
var QueryInfluxDB api.QueryAPI
var DeleteInfluxDB api.DeleteAPI

func ConnectInfluxDB() {
	client := influxdb2.NewClientWithOptions(
		Config().InfluxDBURI,
		Config().InfluxDBToken,
		influxdb2.DefaultOptions().
			SetBatchSize(50).
			SetFlushInterval(200).
			SetUseGZip(true).
			SetRetryInterval(200).
			SetMaxRetries(3).
			SetMaxRetryInterval(500),
	)
	WriteInfluxDB = client.WriteAPI(Config().InfluxDBOrg, Config().InfluxDBDatabase)
	WriteInfluxDBBlocking = client.WriteAPIBlocking(Config().InfluxDBOrg, Config().InfluxDBDatabase)
	QueryInfluxDB = client.QueryAPI(Config().InfluxDBOrg)
	DeleteInfluxDB = client.DeleteAPI()
}
