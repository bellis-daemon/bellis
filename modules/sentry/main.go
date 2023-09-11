package main

import (
	"github.com/bellis-daemon/bellis/common"
	_ "github.com/bellis-daemon/bellis/modules/sentry/apps/implements/all"
	"github.com/bellis-daemon/bellis/modules/sentry/consumer"
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
	consumer.Serve()
}
