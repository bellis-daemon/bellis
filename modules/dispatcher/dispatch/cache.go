package dispatch

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/relock"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/dispatcher/producer"
	"github.com/minoic/glgf"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var mutex *relock.Mutex

const EntityList = "EntityList"

func init() {
	mutex = relock.NewMutex(storage.Redis(), "EntityListMutex")
}

func syncEntityID() {
	ctx := context.Background()
	find, err := storage.CEntity.Find(ctx, bson.M{
		"Active": true,
	})
	if err != nil {
		glgf.Error(err)
		return
	}
	for find.Next(ctx) {
		var entity models.Application
		err := find.Decode(&entity)
		if err != nil {
			glgf.Error(err)
			continue
		}
		func() {
			err = mutex.Lock()
			if err != nil {
				glgf.Error(err)
				return
			}
			defer mutex.Unlock()
			err := storage.Redis().ZRank(ctx, EntityList, entity.ID.Hex()).Err()
			if err != nil {
				glgf.Debug("Entity not fount in redis set EntityList", err)
				err = storage.Redis().ZAdd(ctx, EntityList, redis.Z{
					Score:  timeToScore(time.Now()),
					Member: entity.ID.Hex(),
				}).Err()
				if err != nil {
					glgf.Error(err)
				}
			}
		}()
	}
}

func checkEntities() {
	ctx := context.Background()
	err := mutex.Lock()
	if err != nil {
		glgf.Error(err)
		return
	}
	defer mutex.Unlock()
	result, err := storage.Redis().ZPopMin(ctx, EntityList, 3).Result()
	if err != nil {
		glgf.Error(err)
		return
	}
	for i := range result {
		if scoreToTime(result[i].Score).Before(time.Now()) {
			ddl := time.Now().Add(30 * time.Second)
			result[i].Score = timeToScore(ddl)
			glgf.Debug("entity claiming:", scoreToTime(result[i].Score).Format(time.RFC3339), result[i].Member)
			id, err := primitive.ObjectIDFromHex(cast.ToString(result[i].Member))
			var entity models.Application
			err = storage.CEntity.FindOne(ctx, bson.M{
				"_id": id,
			}).Decode(&entity)
			if err != nil {
				glgf.Error(err)
				continue
			}
			err = producer.EntityClaim(ctx, id.Hex(), ddl, entity)
			if err != nil {
				glgf.Error(err)
			}
		}
		err := storage.Redis().ZAdd(ctx, EntityList, result[i]).Err()
		if err != nil {
			glgf.Error(err)
		}
	}
}

func timeToScore(t time.Time) float64 {
	return cast.ToFloat64(t.UnixMilli())
}

func scoreToTime(s float64) time.Time {
	return time.UnixMilli(cast.ToInt64(s))
}
