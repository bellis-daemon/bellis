package midwares

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func PanicRecover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				debug.PrintStack()
				err = status.Errorf(codes.Internal, "Panic err: %v", e)
			}
		}()
		return handler(ctx, req)
	}
}
