package implements

import (
	mc "github.com/bellis-daemon/bellis/modules/sentry/pkg/minecraft"
	"golang.org/x/net/context"
)

type Minecraft struct {
	Options minecraftOptions
}

func (this *Minecraft) Fetch(ctx context.Context) (any, error) {
	pong, err := mc.Ping(this.Options.Address)
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
	return setOptions(&this.Options)
}

type minecraftStatus struct {
	Version      string `json:"version"`
	OnlinePlayer int    `json:"online_player"`
	MaxPlayer    int    `json:"max_player"`
	Description  string `json:"description"`
	FavIcon      string `json:"fav_icon"`
	ModType      string `json:"mod_type"`
}

type minecraftOptions struct {
	Address string `json:"address"`
}
