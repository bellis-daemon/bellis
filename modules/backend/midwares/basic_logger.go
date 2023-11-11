package midwares

import (
	"context"
	"github.com/minoic/glgf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func BasicLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		start := time.Now()
		resp, err = handler(ctx, req)
		dur := time.Now().Sub(start).Milliseconds()
		var addr string
		if forwarded := md.Get("X-Forwarded-For"); len(forwarded) > 0 {
			addr = forwarded[0]
		} else {
			addr = "Unknown"
		}
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
		md, _ := metadata.FromIncomingContext(ss.Context())
		start := time.Now()
		err := handler(srv, ss)
		dur := time.Now().Sub(start).Milliseconds()
		var addr string
		if forwarded := md.Get("X-Forwarded-For"); len(forwarded) > 0 {
			addr = forwarded[0]
		} else {
			addr = "Unknown"
		}
		if err != nil {
			glgf.Warnf("| %-15s | Streamed <%s> during %d(ms) ERR:%s", addr, info.FullMethod, dur, err.Error())
		} else {
			glgf.Infof("| %-15s | Streamed <%s> during %d(ms)", addr, info.FullMethod, dur)
		}
		return err
	}
}
