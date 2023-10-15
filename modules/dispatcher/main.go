package main

import (
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/dispatcher/consumer"
	"github.com/bellis-daemon/bellis/modules/dispatcher/dispatch"
	"github.com/minoic/glgf"
	_ "net/http/pprof"
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
	go consumer.Serve()
	dispatch.RunTasks()
}
