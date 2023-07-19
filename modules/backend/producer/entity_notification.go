package producer

import (
	"context"
	"encoding/json"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
)

func NoticeEntityUpdate(ctx context.Context, id string, entity *models.Application) error {
	s, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: common.EntityUpdate,
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"id":     id,
			"Entity": s,
		},
	}).Err()
}

func NoticeEntityDelete(ctx context.Context, id string) error {
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: common.EntityDelete,
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"id": id,
		},
	}).Err()
}
