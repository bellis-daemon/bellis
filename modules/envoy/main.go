package main

import (
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/openobserve"
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
	common.AppName = "Envoy"
	glgf.Infof("BuildTime: %s, GoVersion: %s", BuildTime, GoVersion)
	if storage.Config().OpenObserveEnabled {
		openobserve.RegisterGlgf()
	}
}

func main() {
	storage.ConnectMongo()
	storage.ConnectInfluxDB()
	consumer.Serve()
}
