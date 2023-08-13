package main

import (
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/consumer"
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
	consumer.Serve()
}
