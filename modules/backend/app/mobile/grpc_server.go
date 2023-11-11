package mobile

import (
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"github.com/minoic/glgf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

var (
	server *grpc.Server
)

func init() {
	server = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			midwares.PanicRecover(),
			midwares.AuthChecker(),
			midwares.BasicLogger(),
		),
		grpc.ChainStreamInterceptor(
			midwares.AuthCheckerStream(),
			midwares.PanicRecoverStream(),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     10 * time.Second,
			MaxConnectionAge:      15 * time.Second,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:               3 * time.Second,
		}),
	)
}

func Register(reg func(server *grpc.Server) string) {
	glgf.Infof("Registering service to grpc server: %s", reg(server))
}

func ServeGrpc(lis net.Listener) {
	glgf.Success("GRPC server now listening:", lis.Addr())
	err := server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
