package implements

import (
	"github.com/bellis-daemon/bellis/modules/sentry/apps/status"
	mc "github.com/bellis-daemon/bellis/modules/sentry/pkg/minecraft"
	"golang.org/x/net/context"
)

type Minecraft struct {
	options minecraftOptions
}

func (this *Minecraft) Fetch(ctx context.Context) (status.Status, error) {
	pong, err := mc.Ping(this.options.Address)
	if err != nil {
		return &minecraftStatus{}, err
	}
	return &minecraftStatus{
		Version:      pong.Version.Name,
		OnlinePlayer: pong.Players.Online,
		MaxPlayer:    pong.Players.Max,
		Description:  pong.Description.Des,
		FavIcon:      pong.FavIcon,
		ModType:      pong.ModInfo.ModType,
	}, nil
}

func (this *Minecraft) Init(setOptions func(options any) error) error {
	return setOptions(&this.options)
}

type minecraftStatus struct {
	Version      string `json:"version"`
	OnlinePlayer int    `json:"online_player"`
	MaxPlayer    int    `json:"max_player"`
	Description  string `json:"description"`
	FavIcon      string `json:"fav_icon"`
	ModType      string `json:"mod_type"`
}

func (this *minecraftStatus) PullTrigger(triggerName string) *status.TriggerInfo {
	switch triggerName {

	}
	return nil
}

type minecraftOptions struct {
	Address string `json:"address"`
}
