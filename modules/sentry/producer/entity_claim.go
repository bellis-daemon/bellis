package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
	"time"
)

func EntityClaim(ctx context.Context, id string, deadline time.Time, entity *models.Application) error {
	var buf bytes.Buffer
	dec := json.NewEncoder(&buf)
	err := dec.Encode(entity)
	if err != nil {
		return err
	}
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: common.EntityClaim,
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"EntityID": id,
			"Deadline": deadline.Format(time.RFC3339),
			"Entity":   buf.String(),
		},
	}).Err()
}
