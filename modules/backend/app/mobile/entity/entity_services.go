package entity

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/assertion"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
