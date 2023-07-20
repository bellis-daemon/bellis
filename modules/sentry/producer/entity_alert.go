package producer

import (
	"context"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
	"time"
)

func EntityOffline(ctx context.Context, entityID string, message string, offlineTime time.Time) error {
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: common.EntityOfflineAlert,
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"EntityID":    entityID,
			"Message":     message,
			"OfflineTime": offlineTime.UnixMilli(),
		},
	}).Err()
}

func EntityOnline(ctx context.Context, entityID string, onlineTime time.Time) error {
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: common.EntityOnlineAlert,
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"EntityID":    entityID,
			"OfflineTime": onlineTime.UnixMilli(),
		},
	}).Err()
}
