package midwares

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type RecoveryHandlerFunc func(p interface{}) (err error)

// RecoveryHandlerFuncContext is a function that recovers from the panic `p` by returning an `error`.
// The context can be used to extract request scoped metadata and context values.
type RecoveryHandlerFuncContext func(ctx context.Context, p interface{}) (err error)

// PanicRecover returns a new unary server interceptor for panic recovery.
func PanicRecover(opts ...RecoverOption) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				err = recoverFrom(ctx, r, o.recoveryHandlerFunc)
			}
		}()

		resp, err := handler(ctx, req)
		panicked = false
		return resp, err
	}
}

// PanicRecoverStream returns a new streaming server interceptor for panic recovery.
func PanicRecoverStream(opts ...RecoverOption) grpc.StreamServerInterceptor {
	o := evaluateOptions(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				err = recoverFrom(stream.Context(), r, o.recoveryHandlerFunc)
			}
		}()

		err = handler(srv, stream)
		panicked = false
		return err
	}
}

func recoverFrom(ctx context.Context, p interface{}, r RecoveryHandlerFuncContext) error {
	if r == nil {
		return status.Errorf(codes.Internal, "%v", p)
	}
	return r(ctx, p)
}

var (
	defaultOptions = &options{
		recoveryHandlerFunc: func(ctx context.Context, p interface{}) (err error) {
			debug.PrintStack()
			return status.Errorf(codes.Internal, "Panic error: %v", p)
		},
	}
)

type options struct {
	recoveryHandlerFunc RecoveryHandlerFuncContext
}

func evaluateOptions(opts []RecoverOption) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type RecoverOption func(*options)

// WithRecoveryHandler customizes the function for recovering from a panic.
func WithRecoveryHandler(f RecoveryHandlerFunc) RecoverOption {
	return func(o *options) {
		o.recoveryHandlerFunc = RecoveryHandlerFuncContext(func(ctx context.Context, p interface{}) error {
			return f(p)
		})
	}
}

// WithRecoveryHandlerContext customizes the function for recovering from a panic.
func WithRecoveryHandlerContext(f RecoveryHandlerFuncContext) RecoverOption {
	return func(o *options) {
		o.recoveryHandlerFunc = f
	}
}
