package main

import (
	"github.com/bellis-daemon/bellis/modules/sentry/consumer"
)

var (
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	consumer.Serve()
}
