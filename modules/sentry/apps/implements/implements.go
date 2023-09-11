package implements

import (
	"context"
	"fmt"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/minoic/glgf"
	"reflect"
)

// Implement 必须在每个子类中实现的实际方法
type Implement interface {
	// Fetch must return non nil status value, or it will panic
	// return error if entity is offline
	Fetch(ctx context.Context) (status.Status, error)
	Init(setOptions func(options any) error) error
}

var creators = make(map[string]Creator)

type Creator func() Implement

func Add(scheme string, creator Creator) {
	glgf.Infof("adding entity scheme: %s (%s)", scheme, reflect.TypeOf(creator).PkgPath())
	creators[scheme] = creator
}

func Create(ctx context.Context, entity *models.Application) (Implement, error) {
	creator, ok := creators[entity.Scheme]
	if !ok {
		return nil, fmt.Errorf("cant find implement of scheme: %s", entity.Scheme)
	}
	return creator(), nil
}
