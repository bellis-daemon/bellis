package main

import (
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/consumer"
)

var (
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	storage.ConnectMongo()
	consumer.Register()
	redistream.Instance().Serve()
}
