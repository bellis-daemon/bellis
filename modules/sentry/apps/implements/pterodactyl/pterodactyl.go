package pterodactyl

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/minoic/PterodactylGoApi"
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

func (this *Pterodactyl) Init(setOptions func(options any) error) error {
	err := setOptions(&this.options)
	if err != nil {
		return err
	}
	this.client = PterodactylGoApi.NewClient(this.options.Address, this.options.Token)
	return nil
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
	implements.Add("pterodactyl", func() implements.Implement {
		return &Pterodactyl{}
	})
}
