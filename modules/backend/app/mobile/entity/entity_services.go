package entity

import (
	"context"
	"errors"
	"fmt"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/assertion"
	"github.com/bellis-daemon/bellis/modules/backend/producer"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var errEntityOwnership = errors.New("No permission to this entity! ")

func checkEntityOwnershipById(ctx context.Context, user *models.User, entityID string) assertion.AssertionFunc {
	return func() error {
		ok, err := storage.QuickRCSearch(ctx, "EntityOwnership"+user.ID.Hex()+entityID, func() (bool, error) {
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
		})
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
	}
	if src.Public != nil {
		if src.Public.Threshold != nil {
			dst.Public.Threshold = int(*src.Public.Threshold)
		}
		if src.Public.TriggerList != nil {
			dst.Public.TriggerList = src.Public.TriggerList
		}
	}
}

func getEntityUptime(ctx context.Context, entityID string) string {
	s, err := storage.QuickRCSearch(ctx, "Uptime"+entityID, func() (string, error) {
		id, err := primitive.ObjectIDFromHex(entityID)
		if err != nil {
			return cryptoo.FormatDuration(0), err
		}
		var offlineLog models.OfflineLog
		err = storage.COfflineLog.FindOne(ctx, bson.M{"EntityID": id}, options.FindOne().SetSort(bson.M{"_id": -1})).Err()
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				var entity models.Application
				errF := storage.CEntity.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
				if errF != nil {
					return cryptoo.FormatDuration(0), fmt.Errorf("cant find entity by id %s: %w", entityID, errF)
				}
				return cryptoo.FormatDuration(time.Now().Sub(entity.CreatedAt)), nil
			}
			return cryptoo.FormatDuration(0), fmt.Errorf("cant find offline log by EntityID %s: %w", entityID, err)
		}
		if offlineLog.OnlineTime.IsZero() {
			return cryptoo.FormatDuration(0), nil
		}
		return cryptoo.FormatDuration(time.Now().Sub(offlineLog.OnlineTime)), nil
	})
	if err != nil {
		glgf.Error(err)
		return cryptoo.FormatDuration(0)
	}
	return *s
}

func afterDeleteEntity(entityID string) {
	ctx := context.Background()
	id, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		return
	}
	err = storage.DeleteInfluxDB.DeleteWithName(ctx, "bellis", "backend", time.UnixMilli(0), time.Now().Add(time.Hour), fmt.Sprintf(`id="%s"`, entityID))
	if err != nil {
		glgf.Error("error deleting in influxdb", err)
	}
	_, _ = storage.COfflineLog.DeleteMany(ctx, bson.M{"EntityID": id})
	_ = producer.NoticeEntityDelete(ctx, entityID)
}

func afterCreateEntity(entity *models.Application) {
	ctx := context.Background()
	_ = producer.NoticeEntityUpdate(ctx, entity.ID.Hex(), entity)
}

func afterUpdateEntity(entity *models.Application) {
	ctx := context.Background()
	_ = producer.NoticeEntityUpdate(ctx, entity.ID.Hex(), entity)
}
