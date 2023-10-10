package mobile

import (
	"net"
	"time"

	midwares2 "github.com/bellis-daemon/bellis/modules/backend/midwares"
	"github.com/minoic/glgf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	server *grpc.Server
)

func init() {
	server = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ConnectionTimeout(3*time.Second),
		grpc.ChainUnaryInterceptor(
			midwares2.PanicRecover(),
			midwares2.AuthChecker(),
			midwares2.BasicLogger(),
		),
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
