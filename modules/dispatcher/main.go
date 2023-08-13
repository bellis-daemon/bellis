package main

import (
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/dispatcher/consumer"
	"github.com/bellis-daemon/bellis/modules/dispatcher/dispatch"
	_ "net/http/pprof"
)

var (
	BuildTime string
	GoVersion string
)

func init() {
	common.BuildTime = BuildTime
	common.GoVersion = GoVersion
}

func main() {
	storage.ConnectMongo()
	go consumer.Serve()
	dispatch.RunTasks()
}
