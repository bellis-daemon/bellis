package main

import (
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/auth"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/entity"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/profile"
	"github.com/bellis-daemon/bellis/modules/backend/app/server"
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
	server.ServeGrpc()
}
