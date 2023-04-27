package main

import (
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/common/storage"
)

var (
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	storage.ConnectMongo()
	redistream.Instance().Serve()
}
