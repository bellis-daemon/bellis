package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/minoic/glgf"
	"github.com/redis/go-redis/v9"
	"reflect"
	"time"
)

var rdb *redis.Client

func Redis() *redis.Client {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:        "redis:6379",
			DialTimeout: 3 * time.Second,
		})
	}
	return rdb
}

func QuickRCSearch[T any](ctx context.Context, key string, fallback func() (T, error), exp ...time.Duration) (*T, error) {
	const QUICK_RC = "QUICK_RC"
	cmd := Redis().Get(ctx, QUICK_RC+key)
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
			Redis().Set(ctx, QUICK_RC+key, buf.String(), dur)
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
