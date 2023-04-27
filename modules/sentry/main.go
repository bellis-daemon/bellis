package main

import (
	"github.com/bellis-daemon/bellis/common/redistream"
)

var (
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	redistream.Instance().Serve()
}
