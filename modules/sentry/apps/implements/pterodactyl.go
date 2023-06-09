package implements

import (
	"context"
	"github.com/minoic/PterodactylGoApi"
)

type Pterodactyl struct {
	Options pterodactylOptions
	client  *PterodactylGoApi.Client
}

func (this *Pterodactyl) Fetch(ctx context.Context) (any, error) {
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
	err := setOptions(&this.Options)
	if err != nil {
		return err
	}
	this.client = PterodactylGoApi.NewClient(this.Options.Address, this.Options.Token)
	return nil
}

type pterodactylStatus struct {
	UserAmount   int
	ServerAmount int
}

type pterodactylOptions struct {
	Address string `json:"address"`
	Token   string `gorm:"type:blob" json:"token"`
}
