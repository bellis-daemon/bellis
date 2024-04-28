package entity

import (
	"context"
	"errors"
	"fmt"
	"github.com/bellis-daemon/bellis/common/cache"
	"time"

	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/assertion"
	"github.com/bellis-daemon/bellis/modules/backend/producer"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errEntityOwnership = errors.New("no permission to this entity! ")

func checkEntityOwnershipById(ctx context.Context, user *models.User, entityID string) assertion.AssertionFunc {
	return func() error {
		ok, err := storage.QuickRCSearch(ctx, "ENTITY_OWNERSHIP_"+user.ID.Hex()+entityID, func() (bool, error) {
			eid, err := primitive.ObjectIDFromHex(entityID)
			if err != nil {
				return false, err
			}
			var entity models.Application
			err = storage.CEntity.FindOne(ctx, bson.M{"_id": eid}).Decode(&entity)
			if err != nil {
				return false, err
			}
			err = checkEntityOwnership(user, &entity)()
			if err != nil {
				return false, err
			}
			return true, nil
		},time.Hour)
		if err != nil {
			glgf.Warn(err)
			return errEntityOwnership
		}
		if !*ok {
			glgf.Warn("Unauthorized access:", user, entityID)
			return errEntityOwnership
		}
		return nil
	}
}

func checkEntityOwnership(user *models.User, entity *models.Application) assertion.AssertionFunc {
	return func() error {
		if user.ID.Hex() != entity.UserID.Hex() {
			return errEntityOwnership
		}
		return nil
	}
}

func loadPublicOptions(src *Entity, dst *models.Application) {
	dst.Public = models.ApplicationPublicOptions{
		Threshold:   5,
		TriggerList: nil,
		Multiplier:  2,
	}
	if src.Public != nil {
		if src.Public.Threshold != nil {
			dst.Public.Threshold = int(*src.Public.Threshold)
		}
		if src.Public.TriggerList != nil {
			dst.Public.TriggerList = src.Public.TriggerList
		}
		if src.Public.Multiplier != nil {
			dst.Public.Multiplier = uint(*src.Public.Multiplier)
		}
	}
}

func getEntityAvalibility(ctx context.Context, entityID string, duration string) float64 {
	a, err := storage.QuickRCSearch(ctx, "AVALIBILITY_"+entityID, func() (float64, error) {
		result, err := storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(`
total = from(bucket: "backend")
	|> range(start: -%s)
	|> filter(fn: (r) => r["id"] == "%s")
	|> filter(fn: (r) => r["_field"] == "c_live")

total
	|> count()
	|> yield(name: "total")

total
	|> filter(fn: (r) => r["_value"] == true)
	|> count()
	|> yield(name: "live")
		`, duration, entityID))
		if err != nil {
			return 0, err
		}
		var total, live float64 = 0.0, 1.0
		for result.Next() {
			if result.Record().Result() == "total" {
				total = cast.ToFloat64(result.Record().Value())
			} else if result.Record().Result() == "live" {
				live = cast.ToFloat64(result.Record().Value())
			}
		}
		return live / total, nil
	},time.Minute)
	if err != nil {
		glgf.Error(err)
		return 0
	}
	return *a
}

func getEntityUptime(ctx context.Context, entityID string) string {
	s, err := storage.QuickRCSearch(ctx, "UPTIME_"+entityID, func() (string, error) {
		id, err := primitive.ObjectIDFromHex(entityID)
		if err != nil {
			return cryptoo.FormatDuration(0), err
		}
		var offlineLog models.OfflineLog
		err = storage.COfflineLog.FindOne(ctx, bson.M{"EntityID": id}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(&offlineLog)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				var entity models.Application
				errF := storage.CEntity.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
				if errF != nil {
					return cryptoo.FormatDuration(0), fmt.Errorf("cant find entity by id %s: %w", entityID, errF)
				}
				return cryptoo.FormatDuration(time.Since(entity.CreatedAt)), nil
			}
			return cryptoo.FormatDuration(0), fmt.Errorf("cant find offline log by EntityID %s: %w", entityID, err)
		}
		if offlineLog.OnlineTime.IsZero() {
			return cryptoo.FormatDuration(0), nil
		}
		return cryptoo.FormatDuration(time.Since(offlineLog.OnlineTime)), nil
	},time.Minute)
	if err != nil {
		glgf.Error(err)
		return cryptoo.FormatDuration(0)
	}
	return *s
}

func afterDeleteEntity(user *models.User, entityID string) {
	ctx := context.Background()
	id, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		return
	}
	err = storage.DeleteInfluxDB.DeleteWithName(ctx, "bellis", "backend", time.UnixMilli(0), time.Now().Add(time.Hour), fmt.Sprintf(`id="%s"`, entityID))
	if err != nil {
		glgf.Error("error deleting in influxdb", err)
	}
	storage.COfflineLog.DeleteMany(ctx, bson.M{"EntityID": id})
	user.UsageEntityIncr(ctx, -1)
	producer.NoticeEntityDelete(ctx, entityID)
	cache.ExpireUserEntities(ctx, user.ID)
}

func afterCreateEntity(user *models.User, entity *models.Application) {
	ctx := context.Background()
	user.UsageEntityIncr(ctx, 1)
	producer.NoticeEntityUpdate(ctx, entity.ID.Hex(), entity)
	cache.ExpireUserEntities(ctx, user.ID)
}

func afterUpdateEntity(user *models.User, entity *models.Application) {
	ctx := context.Background()
	producer.NoticeEntityUpdate(ctx, entity.ID.Hex(), entity)
	cache.ExpireUserEntities(ctx, user.ID)
}
