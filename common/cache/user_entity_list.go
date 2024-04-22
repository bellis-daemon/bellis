package cache

import (
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"time"
)

const UserEntitiesKey = "UserEntities_"

func GetUserEntities(ctx context.Context, userId primitive.ObjectID) ([]models.Application, error) {
	ret, err := storage.QuickRCSearch[[]models.Application](ctx, UserEntitiesKey+UserEntitiesKey+userId.Hex(), func() ([]models.Application, error) {
		var entities []models.Application
		find, err := storage.CEntity.Find(ctx, bson.M{"UserID": userId})
		if err != nil {
			return nil, err
		}
		err = find.All(ctx, &entities)
		if err != nil {
			return nil, err
		}
		return entities, nil
	}, 10*time.Minute)
	if err != nil {
		return nil, err
	}
	return *ret, nil
}

func ExpireUserEntities(ctx context.Context, userId primitive.ObjectID) error {
	return storage.Redis().Del(ctx, "QUICK_RC_"+UserEntitiesKey+userId.Hex()).Err()
}
