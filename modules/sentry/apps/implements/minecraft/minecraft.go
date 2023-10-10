package minecraft

import (
	"fmt"
	"strings"

	"github.com/bellis-daemon/bellis/modules/sentry/apps/implements"

	"github.com/bellis-daemon/bellis/common"
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
	err := setOptions(&this.options)
	if err != nil {
		return err
	}
	if !strings.Contains(this.options.Address, ":") {
		this.options.Address += ":25565"
	}
	return nil
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
	case "OnlinePlayersThreshold":
		if (float32(this.OnlinePlayer) / float32(this.MaxPlayer)) >= 0.9 {
			return &status.TriggerInfo{
				Name:     "OnlinePlayersThreshold",
				Message:  fmt.Sprintf("Your server`s number of online players exceeds 90%%, now is (%d/%d)", this.OnlinePlayer, this.MaxPlayer),
				Priority: common.Warning,
			}
		}
		return nil
	}
	return nil
}

type minecraftOptions struct {
	Address string `json:"address"`
}

func init() {
	implements.Add("minecraft", func() implements.Implement {
		return &Minecraft{}
	})
}
