package apps

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/bt"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/dns"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/docker"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/http"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/minecraft"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/ping"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/source"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements/vps"
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
		handler = &bt.BT{}
	case Ping:
		handler = &ping.Ping{}
	case HTTP:
		handler = &http.HTTP{}
	case Minecraft:
		handler = &minecraft.Minecraft{}
	case V2Ray:
		handler = &minecraft.Minecraft{}
	case DNS:
		handler = &dns.DNS{}
	case VPS:
		handler = &vps.VPS{}
	case Docker:
		handler = &docker.Docker{}
	case Source:
		handler = &source.Source{}
	}
	return
}
