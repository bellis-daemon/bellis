package main

import (
	"net"

	"github.com/bellis-daemon/bellis/common"
	_ "github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/auth"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/entity"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/profile"
	"github.com/bellis-daemon/bellis/modules/backend/app/web"
	"github.com/minoic/glgf"
	"github.com/soheilhy/cmux"
)

var (
	BuildTime string
	GoVersion string
)

func init() {
	common.BuildTime = BuildTime
	common.GoVersion = GoVersion
	glgf.Infof("BuildTime: %s, GoVersion: %s", BuildTime, GoVersion)
}

func main() {
	storage.ConnectMongo()
	index.InitIndexes()
	storage.ConnectInfluxDB()
	l, err := net.Listen("tcp", "0.0.0.0:7001")
	if err != nil {
		panic(err)
	}
	m := cmux.New(l)
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	webL := m.Match(cmux.HTTP1Fast())
	go mobile.ServeGrpc(grpcL)
	go web.ServeWeb(webL)
	err = m.Serve()
	if err != nil {
		panic(err)
	}
}
