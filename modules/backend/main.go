package main

import (
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"net"

	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/auth"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/entity"
	_ "github.com/bellis-daemon/bellis/modules/backend/app/mobile/profile"
	"github.com/bellis-daemon/bellis/modules/backend/app/web"
	"github.com/soheilhy/cmux"
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
	m.Serve()
}
