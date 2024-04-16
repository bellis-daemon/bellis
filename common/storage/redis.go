package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/minoic/glgf"
	"github.com/redis/go-redis/v9"
)

var rdb redis.UniversalClient

func Redis() redis.UniversalClient {
	if rdb == nil {
		opt := &redis.UniversalOptions{
			Addrs:       Config().RedisAddrs,
			Username:    Config().RedisUsername,
			Password:    Config().RedisPassword,
			DialTimeout: 3 * time.Second,
		}
		if len(opt.Addrs) > 1 {
			opt.ReadOnly = true
			opt.RouteRandomly = true
		}
		rdb = redis.NewUniversalClient(opt)
	}
	return rdb
}

func QuickRCSearch[T any](ctx context.Context, key string, fallback func() (T, error), exp ...time.Duration) (*T, error) {
	const QuickRc = "QUICK_RC_"
	cmd := Redis().Get(ctx, QuickRc+key)
	if cmd.Err() != nil {
		value, err := fallback()
		if err != nil {
			return nil, err
		}
		go func() {
			dur := time.Minute
			if len(exp) != 0 {
				dur = exp[0]
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(map[string]any{
				"t": reflect.TypeOf(value).String(),
				"v": value,
			})
			if err != nil {
				glgf.Warn(err)
				return
			}
			Redis().Set(ctx, QuickRc+key, buf.String(), dur)
		}()
		return &value, nil
	}
	var dec struct {
		Type  string `json:"t"`
		Value T      `json:"v"`
	}
	err := json.Unmarshal([]byte(cmd.Val()), &dec)
	if err != nil {
		return nil, err
	}
	return &dec.Value, nil
}
