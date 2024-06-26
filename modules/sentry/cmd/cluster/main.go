package main

import (
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/openobserve"
	"github.com/bellis-daemon/bellis/common/storage"
	_ "github.com/bellis-daemon/bellis/modules/sentry/apps/implements/all"
	"github.com/bellis-daemon/bellis/modules/sentry/consumer"
	"github.com/minoic/glgf"
)

var (
	BuildTime string
	GoVersion string
)

func init() {
	common.BuildTime = BuildTime
	common.GoVersion = GoVersion
	common.AppName = "Sentry"
	glgf.Infof("BuildTime: %s, GoVersion: %s", BuildTime, GoVersion)
	if storage.Config().OpenObserveEnabled {
		openobserve.RegisterGlgf()
	}
}

func main() {
	storage.ConnectInfluxDB()
	consumer.Serve()
}
