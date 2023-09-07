package factory

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/modules/sentry/apps"
	"sync"
	"time"
)

var entities sync.Map

func RunEntity(entityID string, deadline time.Time, entity *models.Application) error {
	app, err := apps.NewEntity(context.Background(), deadline, entity)
	if err != nil {
		return err
	}
	entities.Store(entityID, app)
	app.Run()
	time.AfterFunc(deadline.Sub(time.Now()), func() {
		entities.Delete(entityID)
	})
	return nil
}

func GetEntity(entityID string) (*apps.Entity, error) {
	entity, ok := entities.Load(entityID)
	if !ok {
		return nil, errors.New("cant find this entity:" + entityID)
	}
	return entity.(*apps.Entity), nil
}

func DeleteEntity(entityID string) {
	entity, ok := entities.Load(entityID)
	if ok {
		entity.(*apps.Entity).Cancel()
		entities.Delete(entityID)
	}
}
