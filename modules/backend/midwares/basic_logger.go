package midwares

import (
	"context"
	"github.com/minoic/glgf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func BasicLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		p, _ := peer.FromContext(ctx)
		resp, err = handler(ctx, req)
		if err != nil {
			glgf.Warnf("| %s |<%s> ERR:%s", p.Addr.String(), info.FullMethod, err.Error())
		} else {
			glgf.Infof("| %s |<%s>", p.Addr.String(), info.FullMethod)
		}
		return
	}
}
