package apps

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
)

// Implement 必须在每个子类中实现的实际方法
type Implement interface {
	// Fetch must return non nil status value, or it will panic
	// return error if entity is offline
	Fetch(ctx context.Context) (status.Status, error)
	Init(setOptions func(options any) error) error
}

func parseImplements(ctx context.Context, entity *models.Application) (handler Implement) {
	switch entity.SchemeID {
	case BT:
		handler = &implements.BT{}
	case Ping:
		handler = &implements.Ping{}
	case HTTP:
		handler = &implements.HTTP{}
	case Minecraft:
		handler = &implements.Minecraft{}
	case V2Ray:
		handler = &implements.Minecraft{}
	case DNS:
		handler = &implements.DNS{}
	case VPS:
		handler = &implements.VPS{}
	case Docker:
		handler = &implements.Docker{}
	case Source:
		handler = &implements.Source{}
	}
	return
}
