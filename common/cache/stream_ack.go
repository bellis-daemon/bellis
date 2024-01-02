package cache

import (
	"context"
	"time"

	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
)

const streamAckPrefix = "STREAM_ACK_"
const streamAckTolerance = 2

type StreamAckEventType uint

const (
	StreamAckSent StreamAckEventType = iota
	StreamAckReceived
)

func StreamAckStart(ctx context.Context, key string) {
	storage.Redis().Set(ctx, streamAckPrefix+key, 0, time.Minute)
}

func StreamAckEvent(ctx context.Context, event StreamAckEventType, key string) {
	switch event {
	case StreamAckReceived:
		storage.Redis().Decr(ctx, streamAckPrefix+key)
	case StreamAckSent:
		storage.Redis().Incr(ctx, streamAckPrefix+key)
	}
	storage.Redis().Expire(ctx, streamAckPrefix+key, time.Minute)
}

func StreamAckStop(ctx context.Context, key string) {
	storage.Redis().Del(ctx, streamAckPrefix+key)
}

func StreamAckCheck(ctx context.Context, key string) bool {
	res, err := storage.Redis().Get(ctx, streamAckPrefix+key).Int()
	if err != nil {
		glgf.Error(err)
		return false
	}
	return res < streamAckTolerance
}
