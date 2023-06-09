package implements

import (
	"context"
	"github.com/rumblefrog/go-a2s"
)

type Source struct {
	Options sourceOptions
	client  *a2s.Client
}

func (this *Source) Fetch(ctx context.Context) (any, error) {
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
	err := setOptions(&this.Options)
	if err != nil {
		return err
	}
	this.client, err = a2s.NewClient(this.Options.Address)
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
