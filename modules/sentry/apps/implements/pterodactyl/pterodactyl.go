package pterodactyl

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/minoic/PterodactylGoApi"
	"go.mongodb.org/mongo-driver/bson"
)

type Pterodactyl struct {
	options pterodactylOptions
	client  *PterodactylGoApi.Client
}

func (this *Pterodactyl) Fetch(ctx context.Context) (status.Status, error) {
	servers, err := this.client.GetAllServers()
	if err != nil {
		return &pterodactylStatus{}, err
	}
	users, err := this.client.GetAllUsers()
	if err != nil {
		return &pterodactylStatus{}, err
	}
	return &pterodactylStatus{
		UserAmount:   len(users),
		ServerAmount: len(servers),
	}, nil
}

type pterodactylStatus struct {
	UserAmount   int
	ServerAmount int
}

func (this *pterodactylStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type pterodactylOptions struct {
	Address string `json:"address"`
	Token   string `gorm:"type:blob" json:"token"`
}

func init() {
	implements.Register("pterodactyl", func(options bson.M) implements.Implement {
		ret := &Pterodactyl{
			options: option.ToOption[pterodactylOptions](options),
			client:  nil,
		}
		ret.client = PterodactylGoApi.NewClient(ret.options.Address, ret.options.Token)
		return ret
	})
}
