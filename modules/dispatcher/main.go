package main

import (
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/openobserve"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/dispatcher/consumer"
	"github.com/bellis-daemon/bellis/modules/dispatcher/dispatch"
	"github.com/minoic/glgf"
)

var (
	BuildTime string
	GoVersion string
)

func init() {
	common.BuildTime = BuildTime
	common.GoVersion = GoVersion
	common.AppName = "Dispatcher"
	glgf.Infof("BuildTime: %s, GoVersion: %s", BuildTime, GoVersion)
	if storage.Config().OpenObserveEnabled {
		openobserve.RegisterGlgf()
	}
}

func main() {
	storage.ConnectMongo()
	go consumer.Serve()
	dispatch.RunTasks()
}
