package main

import (
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/dispatcher/dispatch"
	_ "net/http/pprof"
)

var (
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	storage.ConnectMongo()
	dispatch.RunTasks()
}
