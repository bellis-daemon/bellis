package main

import (
	"github.com/bellis-daemon/bellis/modules/dispatcher/dispatch"
	"net/http"
	_ "net/http/pprof"
)

var (
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	go func() {
		http.ListenAndServe(":6001", nil)

	}()
	dispatch.RunTasks()
}
