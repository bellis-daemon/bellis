package main

import (
	"github.com/bellis-daemon/bellis/common"

	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/consumer"
	"github.com/minoic/glgf"
)

var (
	BuildTime string
	GoVersion string
)

func init() {
	common.BuildTime = BuildTime
	common.GoVersion = GoVersion
	glgf.Infof("BuildTime: %s, GoVersion: %s", BuildTime, GoVersion)
}

func main() {
	storage.ConnectMongo()
	storage.ConnectInfluxDB()
	consumer.Serve()
}
