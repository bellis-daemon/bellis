package producer

import (
	"context"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
	"time"
)

func EntityClaim(ctx context.Context, id string, deadline time.Time) error {
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: "EntityClaim",
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"EntityID": id,
			"Deadline": deadline.Format(time.RFC3339),
		},
	}).Err()
}
