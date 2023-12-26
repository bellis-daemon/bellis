package midwares

import (
	"context"
	"time"

	"github.com/minoic/glgf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func BasicLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		dur := time.Since(start).Milliseconds()
		addr := ipFromContext(ctx)
		if err != nil {
			glgf.Warnf("| %-15s |<%s> in %d(ms) ERR:%s", addr, info.FullMethod, dur, err.Error())
		} else {
			glgf.Infof("| %-15s |<%s> in %d(ms)", addr, info.FullMethod, dur)
		}
		return
	}
}

func BasicLoggerStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		dur := time.Since(start).Milliseconds()
		addr := ipFromContext(ss.Context())
		if err != nil {
			glgf.Warnf("| %-15s | Streamed <%s> during %d(ms) ERR:%s", addr, info.FullMethod, dur, err.Error())
		} else {
			glgf.Infof("| %-15s | Streamed <%s> during %d(ms)", addr, info.FullMethod, dur)
		}
		return err
	}
}

func ipFromContext(ctx context.Context) string {
	var addr string
	md, _ := metadata.FromIncomingContext(ctx)
	if forwarded := md.Get("X-Forwarded-For"); len(forwarded) > 0 {
		addr = forwarded[0]
	} else {
		if p, ok := peer.FromContext(ctx); ok {
			addr = p.Addr.String()
		} else {
			addr = "Unknown"
		}
	}
	return addr
}
