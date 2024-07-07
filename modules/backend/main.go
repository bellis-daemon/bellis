package main

import (
	"context"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/bellis-daemon/bellis/modules/backend/jobs"

	"github.com/bellis-daemon/bellis/common"
	_ "github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/models/index"
	"github.com/bellis-daemon/bellis/common/openobserve"
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
	common.AppName = "Backend"
	glgf.Infof("BuildTime: %s, GoVersion: %s", BuildTime, GoVersion)
	if storage.Config().OpenObserveEnabled {
		openobserve.RegisterGlgf()
	}
}

func main() {
	storage.ConnectMongo()
	index.InitIndexes()
	storage.ConnectInfluxDB()
	jobs.StartAsync()
	l, err := net.Listen("tcp", "0.0.0.0:7002")
	if err != nil {
		panic(err)
	}
	m := cmux.New(l)
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	webL := m.Match(cmux.Any())
	ctx, cancel := context.WithCancel(context.Background())
	go mobile.ServeGrpc(ctx, grpcL)
	go web.ServeWeb(ctx, webL)
	go m.Serve()
	go RunHeadlessChrome()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	glgf.Warn("Shutting down server")
	m.Close()
	cancel()
	time.Sleep(3 * time.Second)
}

func RunHeadlessChrome() {
	cmd := exec.Command(
		"/headless-shell/headless-shell",
		"--no-sandbox",
		"--use-gl=angle",
		"--use-angle=swiftshader",
		"--remote-debugging-address=0.0.0.0",
		"--remote-debugging-port=9222",
	)
	for {
		glgf.Infof("Starting chrome headless shell")
		err := cmd.Run()
		if err != nil {
			glgf.Error("Chrome error: ", err)
		}
		time.Sleep(5 * time.Second)
	}
}
