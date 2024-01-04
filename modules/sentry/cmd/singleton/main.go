package main

import (
	"flag"

	"github.com/bellis-daemon/bellis/common"
	_ "github.com/bellis-daemon/bellis/modules/sentry/apps/implements/all"
	"github.com/bellis-daemon/bellis/modules/sentry/client"
	"github.com/minoic/glgf"
)

var (
	BuildTime string
	GoVersion string

	Token string
)

func init() {
	common.BuildTime = BuildTime
	common.GoVersion = GoVersion
	glgf.Infof("BuildTime: %s, GoVersion: %s", BuildTime, GoVersion)

	flag.StringVar(&Token, "token", "", "Secret token given in bellis app.")
}

func main() {
	flag.Parse()
	if Token == "" {
		panic("token must not be empty")
	}
	client.ServeHttpEventListener(Token)
}
