package implements

import (
	"context"
	"fmt"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

// Implement 必须在每个子类中实现的实际方法
type Implement interface {
	// Fetch must return non nil status value, or it will panic
	// return error if entity is offline
	Fetch(ctx context.Context) (status.Status, error)
	Multiplier() uint
}

var spawners = make(map[string]Spawner)

type Spawner func(options bson.M) Implement

func Register(scheme string, creator Spawner) {
	glgf.Infof("adding entity scheme: %s (%s)", scheme, reflect.TypeOf(creator).PkgPath())
	spawners[scheme] = creator
}

func Spawn(entity *models.Application) (Implement, error) {
	creator, ok := spawners[entity.Scheme]
	if !ok {
		return nil, fmt.Errorf("cant find implement of scheme: %s", entity.Scheme)
	}
	return creator(entity.Options), nil
}

type Template struct{}

func (t Template) Fetch(ctx context.Context) (status.Status, error) {
	panic("Fetch function not implemented")
}

func (t Template) Multiplier() uint {
	return 1
}
