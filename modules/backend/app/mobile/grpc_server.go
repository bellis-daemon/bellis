package mobile

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"github.com/minoic/glgf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"math"
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
			midwares.BasicLoggerStream(),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(math.MaxInt64),
			MaxConnectionAge:      5 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  6 * time.Second,
			Timeout:               3 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: false,
		}),
	)
}

func Register(reg func(server *grpc.Server) string) {
	glgf.Infof("Registering service to grpc server: %s", reg(server))
}

func ServeGrpc(ctx context.Context, lis net.Listener) {
	glgf.Success("GRPC server now listening:", lis.Addr())
	go func() {
		err := server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	select {
	case <-ctx.Done():
		server.GracefulStop()
	}
}

func Server() *grpc.Server {
	return server
}
