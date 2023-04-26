package midwares

import (
	"context"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"time"
)

const (
	CAPACITY           = 1000
	FILL_INTERVAL      = 10 * time.Second / CAPACITY
	WAIT_TIME          = 100 * time.Millisecond
	CHANNEL_LIMIT uint = iota
	TOKENBUCKET_LIMIT
)

var (
	chanContext  = context.Background()
	tokenContext = context.Background()
)

func RateLimiter(mode uint) grpc.UnaryServerInterceptor {
	switch mode {
	case CHANNEL_LIMIT:
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			chans, ok := chanContext.Value(info.Server).(chan struct{})
			if !ok {
				chans = make(chan struct{}, CAPACITY)
				chanContext = context.WithValue(chanContext, info.Server, chans)
			}
			select {
			case chans <- struct{}{}:
				defer func() {
					<-chans
				}()
				return handler(ctx, req)
			case <-time.After(WAIT_TIME):
				return
			}
		}
	case TOKENBUCKET_LIMIT:
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			tokens, ok := tokenContext.Value(info.Server).(*rate.Limiter)
			if !ok {
				tokens = rate.NewLimiter(rate.Every(FILL_INTERVAL), CAPACITY)
				tokenContext = context.WithValue(tokenContext, info.Server, tokens)
			}
			tctx, _ := context.WithTimeout(ctx, WAIT_TIME)
			if tokens.Wait(tctx) == nil {
				return handler(ctx, req)
			}
			return
		}
	}
	panic("invalid limiter mode")
	return nil
}
