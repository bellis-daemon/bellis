package midwares

import (
	"context"
	"github.com/minoic/glgf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"time"
)

func BasicLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		p, _ := peer.FromContext(ctx)
		start := time.Now()
		resp, err = handler(ctx, req)
		dur := time.Now().Sub(start).Milliseconds()
		if err != nil {
			glgf.Warnf("| %s |<%s> in %d(ms) ERR:%s", p.Addr.String(), info.FullMethod, dur, err.Error())
		} else {
			glgf.Infof("| %s |<%s> in %d(ms)", p.Addr.String(), info.FullMethod, dur)
		}
		return
	}
}
