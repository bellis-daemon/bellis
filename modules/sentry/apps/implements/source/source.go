package source

import (
	"context"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/option"
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	"github.com/rumblefrog/go-a2s"
	"go.mongodb.org/mongo-driver/bson"
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

func init() {
	implements.Register("source", func(options bson.M) implements.Implement {
		ret := &Source{options: option.ToOption[sourceOptions](options)}
		var err error
		ret.client, err = a2s.NewClient(ret.options.Address)
		if err != nil {
			panic(err)
		}
		return ret
	})
}
