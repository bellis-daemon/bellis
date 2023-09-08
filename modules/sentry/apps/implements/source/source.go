package source

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/rumblefrog/go-a2s"
)

type Source struct {
	options sourceOptions
	client  *a2s.Client
}

func (this *Source) Fetch(ctx context.Context) (status.Status, error) {
	info, err := this.client.QueryInfo()
	if err != nil {
		return &sourceStatus{}, err
	}
	return &sourceStatus{
		ServerName: info.Name,
		Map:        info.Map,
		Game:       info.Game,
		Players:    int(info.Players),
		MaxPlayers: int(info.MaxPlayers),
		Bots:       int(info.Bots),
		ServerOS:   info.ServerOS.String(),
		Version:    info.Version,
	}, nil
}

func (this *Source) Init(setOptions func(options any) error) error {
	err := setOptions(&this.options)
	if err != nil {
		return err
	}
	this.client, err = a2s.NewClient(this.options.Address)
	return err
}

type sourceOptions struct {
	Address string
}

type sourceStatus struct {
	ServerName string
	Map        string
	Game       string
	Players    int
	MaxPlayers int
	Bots       int
	ServerOS   string
	Version    string
}

func (this *sourceStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}
